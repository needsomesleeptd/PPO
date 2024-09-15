package nn

import "annotater/internal/models"

type INeuralNetwork interface {
	Predict(document models.DocumentData) ([]models.Markup, error)
}
