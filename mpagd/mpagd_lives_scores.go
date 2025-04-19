package mpagd

import (
	"encoding/binary"
	"io"
)

// LivesScore holds the positions of various game elements like score, lives, etc.
type LivesScore struct {
	ScoreTop   uint8
	ScoreLeft  uint8
	LivesTop   uint8
	LivesLeft  uint8
	HighTop    uint8
	HighLeft   uint8
	TimeTop    uint8
	TimeLeft   uint8
	EnergyTop  uint8
	EnergyLeft uint8
}

// LivesScoreInit initializes the LivesScore struct with default values.
// If overwrite is true, it resets the values to defaults.
func (apj *APJFile) LivesScoreInit(overwrite bool) {
	//Implement Livescore
	apj.LivesScore.ScoreTop = 2
	apj.LivesScore.ScoreLeft = 25
	apj.LivesScore.LivesTop = 6
	apj.LivesScore.LivesLeft = 25
	apj.LivesScore.HighTop = 10
	apj.LivesScore.HighLeft = 25
	apj.LivesScore.TimeTop = 25
	apj.LivesScore.TimeLeft = 25
	apj.LivesScore.EnergyTop = 237
	apj.LivesScore.EnergyLeft = 25
}

// readLivesScore reads the LivesScore data from the provided reader.
func (apj *APJFile) readLivesScore(f io.Reader) error {
	// Use a slice of pointers to simplify the reading process.
	fields := []*uint8{
		&apj.LivesScore.ScoreTop, &apj.LivesScore.ScoreLeft,
		&apj.LivesScore.LivesTop, &apj.LivesScore.LivesLeft,
		&apj.LivesScore.HighTop, &apj.LivesScore.HighLeft,
		&apj.LivesScore.TimeTop, &apj.LivesScore.TimeLeft,
		&apj.LivesScore.EnergyTop, &apj.LivesScore.EnergyLeft,
	}

	for _, field := range fields {
		if err := binary.Read(f, binary.LittleEndian, field); err != nil {
			return err
		}
	}
	return nil
}

// writeLivesScore writes the LivesScore data to the provided writer.
func (apj *APJFile) writeLivesScore(f io.Writer) error {
	// Use a slice of values to simplify the writing process.
	fields := []uint8{
		apj.LivesScore.ScoreTop, apj.LivesScore.ScoreLeft,
		apj.LivesScore.LivesTop, apj.LivesScore.LivesLeft,
		apj.LivesScore.HighTop, apj.LivesScore.HighLeft,
		apj.LivesScore.TimeTop, apj.LivesScore.TimeLeft,
		apj.LivesScore.EnergyTop, apj.LivesScore.EnergyLeft,
	}

	for _, field := range fields {
		if err := binary.Write(f, binary.LittleEndian, field); err != nil {
			return err
		}
	}
	return nil
}
