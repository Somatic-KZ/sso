package models

type BankData struct {
	WorkPlace             string `bson:"workPlace" json:"workPlace"`
	CompanyActivityType   string `bson:"companyActivityType" json:"companyActivityType"`
	JobPosition           string `bson:"jobPosition" json:"jobPosition"`
	LastJobWorkExperience string `bson:"lastJobWorkExperience" json:"lastJobWorkExperience"`
	MonthlyIncome         int    `bson:"monthlyIncome" json:"monthlyIncome"`
	AdditionalIncome      int    `bson:"additionalIncome" json:"additionalIncome"`
	ChildrenAmount        uint   `bson:"childrenAmount" json:"childrenAmount"`
	MaritalStatus         string `bson:"maritalStatus" json:"maritalStatus"`
}
