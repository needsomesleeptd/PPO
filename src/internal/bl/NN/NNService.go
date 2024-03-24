package nn

import "annotater/internal/models"

type INeuralNetwork interface {
	Predict(document models.Document) ([]*models.Markup, error)
}
