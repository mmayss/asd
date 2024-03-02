package main

import (
	"fmt"
	"github.com/ericlagergren/decimal"
	"net/http"
)

func main() {
	config := GetConfig()

	if len(config.Repos.Useful) == 0 || len(config.Repos.Useless) == 0 {
		fmt.Println("make sure you have useful and useless repos in your config file")
		return
	}

	client := http.Client{}
	useful := GetAllCommits(client, config.Token, config.Repos.Useful)
	useless := GetAllCommits(client, config.Token, config.Repos.Useless)

	usefulCount, usefulWeightedCount, usefulWeightSum := getCommitCount(useful)
	uselessCount, uselessWeighteCount, uselessWeightSum := getCommitCount(useless)
	sum := usefulCount + uselessCount

	usefulPercentage := new(decimal.Big).SetFloat64((float64(usefulCount) / float64(sum)) * 100).RoundToInt()
	uselessPercentage := new(decimal.Big).SetFloat64((float64(uselessCount) / float64(sum)) * 100).RoundToInt()

	//check if percentages are nan and assign them 0%
	if usefulPercentage.IsNaN(0) {
		usefulPercentage = decimal.New(0, 0)
	}

	if uselessPercentage.IsNaN(0) {
		uselessPercentage = decimal.New(0, 0)
	}

	usefulScore := new(decimal.Big).SetFloat64(float64(usefulWeightedCount) / float64(usefulWeightSum))
	uselessScore := new(decimal.Big).SetFloat64(float64(uselessWeighteCount) / float64(uselessWeightSum))

	fmt.Println("============== counts ===============")
	fmt.Println("useful: ", usefulCount)
	fmt.Println("useless: ", uselessCount)
	fmt.Println("total: ", sum)
	fmt.Println("============ percentage =============")
	fmt.Printf("useful commits: %d %%\n", usefulPercentage)
	fmt.Printf("useless commits: %d %%\n", uselessPercentage)
	fmt.Println("============== score ================")
	fmt.Printf("useful: %d\n", usefulScore)
	fmt.Printf("useless: %d\n", uselessScore)
}

func getCommitCount(repos []Entry) (count int, weightedCount int, weight int) {

	for _, repo := range repos {
		count += repo.CommitCount
		weightedCount += repo.CommitCount * repo.Weight
		weight += repo.Weight
	}

	return count, weightedCount, weight
}
