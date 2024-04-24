package operation

import (
	"alpha-executor/entity"
	"alpha-executor/model"
	"alpha-executor/repository"
	"fmt"
	"github.com/kr/pretty"
	"log"
	"slices"
	"sort"
	"strconv"
)

type Interpreter struct {
	repository *repository.TestingRepository
}

func NewInterpreter(repository *repository.TestingRepository) *Interpreter {
	return &Interpreter{
		repository: repository,
	}
}

func (i *Interpreter) Evaluate(expression Expression) error {
	switch expression.GetKind() {
	case model.PROGRAM.String():
		err := i.evaluateProgram(expression.(*Program))
		if err != nil {
			return err
		}
	default:
		return &entity.CustomError{
			ErrorType: entity.ResponseTypes["CE"],
			Message:   "Program undetected: evaluation failed",
			Position:  entity.Position{},
		}
	}

	return nil
}

func (i *Interpreter) evaluateProgram(program *Program) error {
	for _, expression := range program.body {
		_, err := i.evaluateExpression(expression)
		if err != nil {
			return err
		}
	}

	if len(program.body) == 0 {
		return &entity.CustomError{
			ErrorType: entity.ResponseTypes["CE"],
			Message:   fmt.Sprint("No query to compile"),
			Position:  entity.Position{},
		}
	}

	return nil
}

func (i *Interpreter) evaluateExpression(expression Expression) (bool, error) {
	switch expression.GetKind() {
	case model.GET.String(), model.HOLD.String():
		return i.evaluateGet(expression.(*GetHoldExpression), expression.GetKind())
	case model.EQUALS.String(),
		model.NOT_EQUALS.String(),
		model.LESS_THAN_EQUALS.String(),
		model.GREATER_THAN_EQUALS.String(),
		model.LESS_THAN.String(),
		model.GREATER_THAN.String():
		return i.evaluateComparison(expression.(*BinaryExpression))
	case model.RANGE.String():
		return i.evaluateRange(expression.(*RangeExpression))
	case model.CONJUNCTION.String():
		return i.evaluateConjunction(expression.(*BinaryExpression))
	case model.DISJUNCTION.String():
		return i.evaluateDisjunction(expression.(*BinaryExpression))
	case model.EXISTS.String():
		return i.evaluateExists(expression.(*BinaryExpression))
	case model.FOR_ALL.String():
		return i.evaluateForAll(expression.(*BinaryExpression))
	case model.NEGATION.String():
		return i.evaluateNegation(expression.(*UnaryExpression))
	case model.IMPLICATION.String():
		return i.evaluateImplication(expression.(*BinaryExpression))
	case model.ASSIGN.String():
		return i.evaluateAssignment(expression.(*BinaryExpression))
	case model.UPDATE.String():
		return i.evaluateUpdate(expression.(*UnaryExpression))
	case model.RELEASE.String():
		return i.evaluateRelease(expression.(*UnaryExpression))
	default:
		return false, &entity.CustomError{
			ErrorType: entity.ResponseTypes["RT"],
			Message:   fmt.Sprintf("Unknown kind %s", expression.GetKind()),
		}
	}
}

func (i *Interpreter) evaluateFreeRelation(
	relations []string,
	expression Expression,
	resultRelations *entity.Relations,
) (bool, error) {
	relationName := relations[0]
	relations = relations[1:]
	relation, err := i.repository.GetRelation(relationName)
	if err != nil {
		return false, err
	}

	newRelation := make(entity.Relation)
	for row := range *relation {
		result := false
		rowCopy := *row
		i.repository.AddRow(relationName, &rowCopy)

		if len(relations) == 0 {
			result, err = i.evaluateExpression(expression)
			if err != nil {
				return false, err
			}
		}

		if len(relations) > 0 {
			result, err = i.evaluateFreeRelation(relations, expression, resultRelations)
			if err != nil {
				return false, err
			}
		}

		if result {
			newRelation[&rowCopy] = struct{}{}
		}
	}

	if _, exists := (*resultRelations)[relationName]; exists {
		currentRelation := (*resultRelations)[relationName]
		for row := range newRelation {
			(*currentRelation)[row] = struct{}{}
		}
	} else {
		(*resultRelations)[relationName] = &newRelation
	}

	if len(newRelation) > 0 {
		return true, nil
	}

	return false, nil
}

func (i *Interpreter) evaluateGet(expression *GetHoldExpression, operation string) (bool, error) {
	relation := expression.variable.(*IdentifierExpression)
	if relation.kind != model.FREE_RELATION.String() {
		return false, &entity.CustomError{
			ErrorType: entity.ResponseTypes["RT"],
			Message:   fmt.Sprintf("Expected relation"),
			Position:  relation.position,
		}
	}

	relations, attributes, isRelation, err := i.getData(expression, relation)
	if err != nil {
		return false, err
	}

	resultRelations := make(entity.Relations)
	evaluationResult, err := i.evaluateFreeRelation(relations, expression.expression, &resultRelations)
	if err != nil {
		err.(*entity.CustomError).Position = expression.position
		return false, err
	}

	i.repository.AddCalculatedRelations(resultRelations)

	sortExpression := expression.sort.(*UnaryExpression)
	sorted := make([]*entity.RowMap, 0)

	var result *entity.Relation
	if !isRelation && evaluationResult {
		rel, err := i.joiningRelations(relations)
		if err != nil {
			return false, err
		}

		sorted, err = i.evaluateSort(sortExpression, rel)
		if err != nil {
			return false, err
		}

		projection := Projection{}
		relationPair := entity.Pair[string, *entity.Relation]{Left: relation.value, Right: rel}
		result, err = projection.Execute(relationPair, attributes, relation.position)
		if err != nil {
			return false, err
		}
	} else if isRelation && evaluationResult {
		result, err = i.repository.GetRelation(relations[0])
		if err != nil {
			return false, err
		}

		sorted, err = i.evaluateSort(sortExpression, result)
		if err != nil {
			return false, err
		}
	} else {
		if err = i.addToRepository(expression, operation, relation.value, result); err != nil {
			return false, err
		}

		return true, nil
	}

	if _, err = pretty.Print(sorted); err != nil {
		return false, err
	}

	resultRowNum := expression.rows.(*IdentifierExpression).value
	if err = i.limitResultRows(result, resultRowNum); err != nil {
		return false, err
	}

	if err = i.addToRepository(expression, operation, relation.value, result); err != nil {
		return false, err
	}

	return true, nil
}

func (i *Interpreter) addToRepository(
	expression *GetHoldExpression,
	operation string,
	relationName string,
	result *entity.Relation,
) error {
	if result == nil {
		return nil
	}

	switch operation {
	case model.GET.String():
		i.repository.AddGetRelation(relationName, result)
		break
	case model.HOLD.String():
		i.repository.AddHeldRelation(relationName, result)
	default:
		return &entity.CustomError{
			ErrorType: entity.ResponseTypes["RT"],
			Message:   fmt.Sprintf("Unsupported operation %s", operation),
			Position:  expression.position,
		}
	}
	return nil
}

func (i *Interpreter) getData(expression *GetHoldExpression, relation *IdentifierExpression) ([]string, []string, bool, error) {
	relations := make([]string, 0)
	attributes := make([]string, 0)
	isRelation := false
	for _, data := range expression.relations {
		switch data.GetKind() {
		case model.FREE_RELATION.String():
			isRelation = true
			relations = append(relations, data.(*IdentifierExpression).value)
			break
		case model.ATTRIBUTE.String():
			attr := model.Attribute{}
			assertedData := data.(*IdentifierExpression)
			attribute, err := attr.ExtractAttribute(assertedData.value, assertedData.position)
			if err != nil {
				err.(*entity.CustomError).Position = data.(*IdentifierExpression).position
				return nil, nil, false, err
			}

			if !slices.Contains(relations, attribute.Relation) {
				relations = append(relations, attribute.Relation)
			}

			if !slices.Contains(relations, attribute.Attribute) {
				attributes = append(attributes, attribute.Attribute)
			}
			break
		default:
			return nil, nil, false, &entity.CustomError{
				ErrorType: entity.ResponseTypes["CE"],
				Message:   "Unexpected type",
				Position:  relation.position,
			}
		}

		if isRelation {
			break
		}
	}

	return relations, attributes, isRelation, nil
}

func (i *Interpreter) joiningRelations(relations []string) (*entity.Relation, error) {
	join := Join{}
	for _, rel1Name := range relations {
		rel1Value, err := i.repository.GetCalculatedRelation(rel1Name)
		if err != nil {
			return nil, err
		}

		rel1Pair := entity.Pair[string, *entity.Relation]{Left: rel1Name, Right: rel1Value}
		relations = relations[1:]
		for _, rel2Name := range relations {
			rel2Value, err := i.repository.GetCalculatedRelation(rel2Name)
			if err != nil {
				return nil, err
			}

			rel2Pair := entity.Pair[string, *entity.Relation]{Left: rel2Name, Right: rel2Value}
			relations = relations[1:]
			commonAttributes := make([]string, 0)
			for row1 := range *rel1Pair.Right {
				for row2 := range *rel2Pair.Right {
					exists := make(map[string]struct{})
					for key := range *row1 {
						exists[key] = struct{}{}
					}

					for key := range *row2 {
						if _, ok := exists[key]; ok {
							commonAttributes = append(commonAttributes, key)
						}
					}
					break
				}
				break
			}

			rel1Pair.Right, err = join.Execute(rel1Pair, rel2Pair, commonAttributes)
			if err != nil {
				log.Fatal(err)
			}
		}
		return rel1Pair.Right, nil
	}

	return nil, &entity.CustomError{
		ErrorType: entity.ResponseTypes["RT"],
		Message:   "Relation join error",
		Position:  entity.Position{},
	}
}

func (i *Interpreter) mapToSlice(relation *entity.Relation) []*entity.RowMap {
	relationSlice := make([]*entity.RowMap, 0, len(*relation))
	for row := range *relation {
		relationSlice = append(relationSlice, row)
	}

	return relationSlice
}

func (i *Interpreter) limitResultRows(result *entity.Relation, resultRowNum string) error {
	resultSliced := make(entity.Relation)
	if resultRowNum != model.NULL.String() {
		rowNum, err := strconv.Atoi(resultRowNum)
		if err != nil {
			return err
		}

		if rowNum < len(*result) {
			count := 0
			for row := range *result {
				if count >= rowNum {
					break
				}

				resultSliced[row] = struct{}{}
				count++
			}

			result = &resultSliced
		}
	}

	return nil
}

func (i *Interpreter) evaluateComparison(expression *BinaryExpression) (bool, error) {
	comparison := NewComparison(i.repository)
	return comparison.Compare(expression)
}

func (i *Interpreter) evaluateRange(expression *RangeExpression) (bool, error) {
	relation, err := i.repository.GetRelation(expression.relation.(*IdentifierExpression).value)
	if err != nil {
		return false, err
	}

	i.repository.AddRelation(expression.variable.(*IdentifierExpression).value, relation)
	return true, nil
}

func (i *Interpreter) evaluateConjunction(expression *BinaryExpression) (bool, error) {
	left, err := i.evaluateExpression(expression.left)
	if err != nil {
		return false, err
	}

	right, err := i.evaluateExpression(expression.right)
	if err != nil {
		return false, err
	}

	return left && right, nil
}

func (i *Interpreter) evaluateDisjunction(expression *BinaryExpression) (bool, error) {
	left, err := i.evaluateExpression(expression.left)
	if err != nil {
		return false, err
	}

	right, err := i.evaluateExpression(expression.right)
	if err != nil {
		return false, err
	}

	return left || right, nil
}

func (i *Interpreter) evaluateExists(expression *BinaryExpression) (bool, error) {
	relationName := expression.left.(*IdentifierExpression).value
	left, err := i.repository.GetRelation(relationName)
	if err != nil {
		return false, err
	}

	for row := range *left {
		i.repository.AddRow(relationName, row)

		right, err := i.evaluateExpression(expression.right)
		if err != nil {
			return false, err
		}

		if right {
			return true, nil
		}
	}

	return false, nil
}

func (i *Interpreter) evaluateForAll(expression *BinaryExpression) (bool, error) {
	relationName := expression.left.(*IdentifierExpression).value
	left, err := i.repository.GetRelation(relationName)
	if err != nil {
		return false, err
	}

	for row := range *left {
		i.repository.AddRow(relationName, row)

		right, err := i.evaluateExpression(expression.right)
		if err != nil {
			return false, err
		}

		if !right {
			return false, nil
		}
	}

	return true, nil
}

func (i *Interpreter) evaluateNegation(expression *UnaryExpression) (bool, error) {
	calculatedExpression, err := i.evaluateExpression(expression.expression)
	if err != nil {
		return false, err
	}

	return !calculatedExpression, nil
}

func (i *Interpreter) evaluateImplication(expression *BinaryExpression) (bool, error) {
	left, err := i.evaluateExpression(expression.left)
	if err != nil {
		return false, err
	}

	right, err := i.evaluateExpression(expression.right)
	if err != nil {
		return false, err
	}

	return !left || right, nil
}

func (i *Interpreter) evaluateSort(expression *UnaryExpression, relation *entity.Relation) ([]*entity.RowMap, error) {
	if expression.kind == model.NULL.String() {
		return nil, nil
	}

	if expression.kind != model.UP.String() && expression.kind != model.DOWN.String() {
		return nil, &entity.CustomError{
			ErrorType: entity.ResponseTypes["CE"],
			Message:   "Such sort type is undefined",
			Position:  expression.position,
		}
	}

	attr := model.Attribute{}
	complexAttribute, err := attr.ExtractAttribute(expression.expression.(*IdentifierExpression).value, expression.position)
	if err != nil {
		return nil, err
	}

	attribute := complexAttribute.Attribute
	for row := range *relation {
		if _, exists := (*row)[attribute]; !exists {
			return nil, &entity.CustomError{
				ErrorType: entity.ResponseTypes["CE"],
				Message:   fmt.Sprintf("Attribute %s doesn't exist", attribute),
			}
		}
	}

	relationSlice := i.mapToSlice(relation)
	sort.Slice(relationSlice, func(i, j int) bool {
		values1, values2 := (*relationSlice[i])[attribute], (*relationSlice[j])[attribute]
		for _, value1 := range values1 {
			for _, value2 := range values2 {
				switch expression.kind {
				case model.UP.String():
					if value1 > value2 {
						return true
					} else {
						return false
					}
				case model.DOWN.String():
					if value1 < value2 {
						return true
					} else {
						return false
					}
				}
			}
		}

		return false
	})

	return relationSlice, nil
}

func (i *Interpreter) evaluateAssignment(expression *BinaryExpression) (bool, error) {
	relationAttribute := expression.left.(*IdentifierExpression).value
	assignedValue := expression.right.(*IdentifierExpression).value

	attr := model.Attribute{}
	complexAttribute, err := attr.ExtractAttribute(relationAttribute, expression.position)
	if err != nil {
		return false, err
	}

	relation, err := i.repository.GetHeldRelation(complexAttribute.Relation)
	if err != nil {
		return false, err
	}

	for row := range *relation {
		if _, exists := (*row)[complexAttribute.Attribute]; !exists {
			return false, &entity.CustomError{
				ErrorType: entity.ResponseTypes["CE"],
				Message:   fmt.Sprintf("Attribute %s doesn't exist", complexAttribute.Attribute),
				Position:  expression.position,
			}
		}

		(*row)[complexAttribute.Attribute] = []string{assignedValue}
	}

	return true, nil
}

func (i *Interpreter) evaluateUpdate(expression *UnaryExpression) (bool, error) {
	relationName := expression.expression.(*IdentifierExpression).value
	relation, err := i.repository.GetHeldRelation(relationName)
	if err != nil {
		return false, err
	}

	i.repository.AddRelation(relationName, relation)
	return true, nil
}

func (i *Interpreter) evaluateRelease(expression *UnaryExpression) (bool, error) {
	relationName := expression.expression.(*IdentifierExpression).value
	if _, err := i.repository.GetHeldRelation(relationName); err != nil {
		return false, err
	}

	i.repository.ReleaseHeldRelation(relationName)
	return true, nil
}
