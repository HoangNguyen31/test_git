package main

import (
    "fmt"
    "math/rand"
    "os"
    "path/filepath"
    "time"

    "github.com/google/uuid"

    "github.com/hyperjumptech/grule-rule-engine/ast"
    "github.com/hyperjumptech/grule-rule-engine/builder"
    "github.com/hyperjumptech/grule-rule-engine/engine"
    "github.com/hyperjumptech/grule-rule-engine/pkg"
)

type User struct {
    ID         string
    Name       string
    Gender     string
    Membership string
    Discount   string
}

func generateRandomUsers(count int, nameList, genderList, membershipList []string) []*User {
    rand.Seed(time.Now().UnixNano())

    userList := make([]*User, count)
    for i := 0; i < count; i++ {
        userList[i] = &User{
            ID:         randomID(),
            Name:       randomName(nameList),
            Gender:     randomGender(genderList),
            Membership: randomMembership(membershipList),
            Discount:   "",
        }
    }
    return userList
}

func randomID() string {
    uid, _ := uuid.NewRandom()
    return uid.String()
}

func randomName(seedList []string) string {
    return seedList[rand.Intn(len(seedList))]
}

func randomGender(seedList []string) string {
    return seedList[rand.Intn(len(seedList))]
}

func randomMembership(seedList []string) string {
    return seedList[rand.Intn(len(seedList))]
}

func printStatistics(userList []*User) {
    womanDiscountCount := 0
    goldDiscountCount := 0
    silverDiscountCount := 0
    bronzeDiscountCount := 0
    for _, user := range userList {
        if user.Gender == "WOMAN" && user.Discount != "" {
            womanDiscountCount++
        }
        if user.Membership == "GOLD" && user.Discount != "" && user.Gender != "WOMAN" {
            goldDiscountCount++
        }
        if user.Membership == "SILVER" && user.Discount != "" && user.Gender != "WOMAN" {
            silverDiscountCount++
        }
        if user.Membership == "BRONZE" && user.Discount != "" && user.Gender != "WOMAN" {
            bronzeDiscountCount++
        }
    }
    fmt.Printf("Number of women with discount: %d\n", womanDiscountCount)
    fmt.Printf("Number of gold members with discount: %d\n", goldDiscountCount)
    fmt.Printf("Number of silver members with discount: %d\n", silverDiscountCount)
    fmt.Printf("Number of bronze members with discount: %d\n", bronzeDiscountCount)
	fmt.Printf("Total: %d\n", womanDiscountCount + goldDiscountCount + silverDiscountCount + bronzeDiscountCount)
}

const (
    ruleName    = "Calculate Discount"
    ruleVersion = "0.0.2"
    rulePath    = "./rule/rule.grl"
)

func main() {
    nameList := []string{"Alice", "Bob", "Charlie", "David", "Emma", "Frank"}
    genderList := []string{"MAN", "WOMAN"}
    membershipList := []string{"GOLD", "SILVER", "BRONZE"}

    randomUsers := generateRandomUsers(100000, nameList, genderList, membershipList)

    knowledgeLibrary := ast.NewKnowledgeLibrary()
    ruleBuilder := builder.NewRuleBuilder(knowledgeLibrary)

    path, _ := filepath.Abs(rulePath)
    f, _ := os.Open(path)
    rs := pkg.NewReaderResource(f)

    _ = ruleBuilder.BuildRuleFromResource(ruleName, ruleVersion, rs)

    knowledgeBase, _ := knowledgeLibrary.NewKnowledgeBaseInstance(ruleName, ruleVersion)
    ruleEngine := engine.NewGruleEngine()

    for _, user := range randomUsers {
        dataCtx := ast.NewDataContext()
        _ = dataCtx.Add("User", user)
        err := ruleEngine.Execute(dataCtx, knowledgeBase)
        if err != nil {
            fmt.Println("Error executing rules:", err)
        }
    }

    printStatistics(randomUsers)
}
