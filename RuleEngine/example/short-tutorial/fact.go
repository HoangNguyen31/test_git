package main

import (
	"fmt"
	"time"

	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/engine"
	"github.com/hyperjumptech/grule-rule-engine/pkg"
)

type MyFact struct {
	IntAttribute     int64
	FloatAttribute   float64
	StringAttribute  string
	BooleanAttribute bool
	TimeAttribute    time.Time
	WhatToSay        string
}

func (mf *MyFact) GetWhatToSay(sentence string) string {
	return fmt.Sprintf("Let say \"%s\"", sentence)
}

func main() {
	// Add Fact Into DataContext
	myFact := &MyFact{
		IntAttribute:     123,
		StringAttribute:  "Some string value",
		BooleanAttribute: true,
		FloatAttribute:   1.234,
		TimeAttribute:    time.Now(),
		WhatToSay:        "Hello World",
	}

	dataCtx := ast.NewDataContext()
	err := dataCtx.Add("MF", myFact)
	if err != nil {
		panic(err)
	}

	// Creating a KnowledgeLibrary and Adding Rules Into It
	knowledgeLibrary := ast.NewKnowledgeLibrary()
	ruleBuilder := builder.NewRuleBuilder(knowledgeLibrary)

	rule_1 := `
	rule CheckValuesForHello "Check the default values for hello" salience 10 {
    when 
        MF.IntAttribute == 123 && MF.StringAttribute == "Some string value"
    then
        MF.WhatToSay = MF.GetWhatToSay("Hello Grule");
        Retract("CheckValuesForHello");
	}
	`

	rule_2 := `
	rule CheckValuesForBye "Check the default values for bye" salience 10 {
    when 
        MF.FloatAttribute == 20.11 && MF.StringAttribute == "Some string value"
    then
        MF.WhatToSay = MF.GetWhatToSay("Bye Grule");
        Retract("CheckValuesForBye");
	}
	`
	bs := pkg.NewBytesResource([]byte(rule_1 + rule_2)) // From String or ByteArray
	err = ruleBuilder.BuildRuleFromResource("TutorialRules", "0.0.1", bs)
	if err != nil {
		panic(err)
	}

	// Print rules in the knowledgeLibrary
	rules := knowledgeLibrary.GetKnowledgeBase("TutorialRules", "0.0.1").RuleEntries
	fmt.Println("Number of rules:", len(rules))
	fmt.Println("---------------------------")
	for _, rule := range rules {
		fmt.Printf("Rule Name: %v\n", rule.RuleName)
		fmt.Printf("Rule Description: %v\n", rule.RuleDescription)
		fmt.Printf("Rule Salience: %v\n", rule.Salience)
		fmt.Printf("Rule GrlText: %v\n", rule.GrlText)
		fmt.Println("---------------------------")
	}

	// Executing Grule Rule Engine
	knowledgeBase, _ := knowledgeLibrary.NewKnowledgeBaseInstance("TutorialRules", "0.0.1")

	egn := engine.NewGruleEngine()
	err = egn.Execute(dataCtx, knowledgeBase)
	if err != nil {
		panic(err)
	}

	// Obtaining Result 1
	fmt.Println(myFact)

	// Obtaining Rule 2
	myFact.FloatAttribute = 20.11
	err = egn.Execute(dataCtx, knowledgeBase)
	if err != nil {
		panic(err)
	}
	fmt.Println(myFact)
	fmt.Println()

	// Resources
	/// From File
	fileRes := pkg.NewFileResource("./rule/rules.grl")
	err = ruleBuilder.BuildRuleFromResource("TutorialRules", "0.0.1", fileRes)
	if err != nil {
		panic(err)
	}
	rules = knowledgeLibrary.GetKnowledgeBase("TutorialRules", "0.0.1").RuleEntries
	fmt.Println("Number of rules after getting from files:", len(rules))

	// From URL

	// From GIT
	bundle := pkg.NewGITResourceBundle("https://github.com/hyperjumptech/grule-rule-engine.git", "/examples/benchmark/100_rules.grl")
	resources := bundle.MustLoad()
	for _, res := range resources {
		err := ruleBuilder.BuildRuleFromResource("TutorialRules", "0.0.1", res)
		if err != nil {
			panic(err)
		}
	}
	rules = knowledgeLibrary.GetKnowledgeBase("TutorialRules", "0.0.1").RuleEntries
	fmt.Println("Number of rules after getting from git:", len(rules))
}
