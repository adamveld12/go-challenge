package drum

import (
	"fmt"
	"io/ioutil"
)

// DecodeFile decodes the drum machine file found at the provided path
// and returns a pointer to a parsed pattern which is the entry point to the
// rest of the data.
// TODO: implement
func DecodeFile(path string) (*Pattern, error) {
	data, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}

	p := &Pattern{}
	parseErr := parse(data, p)

	if parseErr != nil {
		return nil, parseErr
	}

	return p, nil
}

// Pattern is the high level representation of the
// drum pattern contained in a .splice file.
type Pattern struct {
	Version  string
	Tempo    float32
	Channels []Channel
}

func (p *Pattern) String() string {
	output := ""

	output += fmt.Sprintf("Saved with HW Version: %s\n", p.Version)
	output += fmt.Sprint("Tempo: ", p.Tempo, "\n")

	for _, channel := range p.Channels {
		output += fmt.Sprintf("%s\n", channel.String())
	}

	return output
}

type Channel struct {
	Id    int32
	Name  string
	Steps []uint32
}

func (c *Channel) printPattern() string {
	output := "|"

	// for each int32
	for _, val := range c.Steps {
		mask := uint32(0x000000FF)

		// for each byte
		for idx := uint32(0); idx < 4; idx++ {
			var shiftedMask uint32 = mask << (idx * 8)
			if val&shiftedMask == 0x00 {
				output += "-"
			} else {
				output += "x"
			}
		}

		output += "|"
	}

	return output
}

func (c *Channel) String() string {
	return fmt.Sprintf("(%d) %s\t%s", c.Id, c.Name, c.printPattern())
}
