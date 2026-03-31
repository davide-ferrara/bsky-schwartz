package models

import "fmt"

type SchwartzWeights struct {
	Reputation          float64 `json:"Reputation"`
	Power               float64 `json:"Power"`
	Wealth              float64 `json:"Wealth"`
	Achievement         float64 `json:"Achievement"`
	Pleasure            float64 `json:"Pleasure"`
	IndependentThoughts float64 `json:"Independent thoughts"`
	IndependentActions  float64 `json:"Independent actions"`
	Stimulation         float64 `json:"Stimulation"`
	PersonalSecurity    float64 `json:"Personal security"`
	SocietalSecurity    float64 `json:"Societal security"`
	Tradition           float64 `json:"Tradition"`
	Lawfulness          float64 `json:"Lawfulness"`
	Respect             float64 `json:"Respect"`
	Humility            float64 `json:"Humility"`
	Responsibility      float64 `json:"Responsibility"`
	Caring              float64 `json:"Caring"`
	Equality            float64 `json:"Equality"`
	Nature              float64 `json:"Nature"`
	Tolerance           float64 `json:"Tolerance"`
}

func MapFormToWeights(formData map[string]string) SchwartzWeights {
	return SchwartzWeights{
		Reputation:          parseFloat(formData["face"]),
		Power:               parseFloat(formData["dominance"]),
		Wealth:              parseFloat(formData["resources"]),
		Achievement:         parseFloat(formData["achievement"]),
		Pleasure:            parseFloat(formData["hedonism"]),
		IndependentThoughts: parseFloat(formData["sd_thought"]),
		IndependentActions:  parseFloat(formData["sd_action"]),
		Stimulation:         parseFloat(formData["stimulation"]),
		PersonalSecurity:    parseFloat(formData["personal_sec"]),
		SocietalSecurity:    parseFloat(formData["societal_sec"]),
		Tradition:           parseFloat(formData["tradition"]),
		Lawfulness:          parseFloat(formData["rule_conf"]),
		Respect:             parseFloat(formData["inter_conf"]),
		Humility:            parseFloat(formData["humility"]),
		Responsibility:      parseFloat(formData["dependability"]),
		Caring:              parseFloat(formData["caring"]),
		Equality:            parseFloat(formData["universalism"]),
		Nature:              parseFloat(formData["nature"]),
		Tolerance:           parseFloat(formData["tolerance"]),
	}
}

func parseFloat(s string) float64 {
	var f float64
	_, err := fmt.Sscanf(s, "%f", &f)
	if err != nil {
		return 0
	}
	return f
}
