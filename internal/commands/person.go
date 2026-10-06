package commands

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"strings"
	"unicode"
)

// PersonCommand ------------------------------------------------------------------------

type PersonResult struct {
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Gender       string `json:"gender"`
	DateOfBirth  string `json:"date_of_birth"`
	PlaceOfBirth string `json:"place_of_birth"`
	FiscalCode   string `json:"fiscal_code"`
}

func (r PersonResult) String() string {
	var sb strings.Builder
	_, _ = fmt.Fprintf(&sb, "First Name:     %s\n", r.FirstName)
	_, _ = fmt.Fprintf(&sb, "Last Name:      %s\n", r.LastName)
	_, _ = fmt.Fprintf(&sb, "Gender:         %s\n", r.Gender)
	_, _ = fmt.Fprintf(&sb, "Date of Birth:  %s\n", r.DateOfBirth)
	_, _ = fmt.Fprintf(&sb, "Place of Birth: %s\n", r.PlaceOfBirth)
	_, _ = fmt.Fprintf(&sb, "Fiscal Code:    %s", r.FiscalCode)
	return sb.String()
}

type PersonCommand struct{}

func (cmd PersonCommand) Execute(otherArgs map[string]string) (fmt.Stringer, error) {
	country := otherArgs["country"]
	if country != "IT" {
		return nil, errors.New("Invalid country: " + country + ". Only IT is supported at the moment")
	}

	result := generateItalianPerson(rand.Intn)
	return result, nil
}

func (cmd PersonCommand) DiscardOutput() bool {
	return false
}

type NameGender struct {
	Name   string
	Gender string
}

type PlaceCode struct {
	Name string
	Code string
}

var italianNames = []NameGender{
	{Name: "Marco", Gender: "M"},
	{Name: "Giuseppe", Gender: "M"},
	{Name: "Francesco", Gender: "M"},
	{Name: "Alessandro", Gender: "M"},
	{Name: "Andrea", Gender: "M"},
	{Name: "Lorenzo", Gender: "M"},
	{Name: "Matteo", Gender: "M"},
	{Name: "Gabriele", Gender: "M"},
	{Name: "Mattia", Gender: "M"},
	{Name: "Leonardo", Gender: "M"},
	{Name: "Davide", Gender: "M"},
	{Name: "Riccardo", Gender: "M"},
	{Name: "Federico", Gender: "M"},
	{Name: "Luca", Gender: "M"},
	{Name: "Stefano", Gender: "M"},
	{Name: "Antonio", Gender: "M"},
	{Name: "Roberto", Gender: "M"},
	{Name: "Mario", Gender: "M"},
	{Name: "Giovanni", Gender: "M"},
	{Name: "Paolo", Gender: "M"},
	{Name: "Anna", Gender: "F"},
	{Name: "Maria", Gender: "F"},
	{Name: "Giulia", Gender: "F"},
	{Name: "Sofia", Gender: "F"},
	{Name: "Aurora", Gender: "F"},
	{Name: "Alice", Gender: "F"},
	{Name: "Ginevra", Gender: "F"},
	{Name: "Emma", Gender: "F"},
	{Name: "Giorgia", Gender: "F"},
	{Name: "Chiara", Gender: "F"},
	{Name: "Sara", Gender: "F"},
	{Name: "Martina", Gender: "F"},
	{Name: "Laura", Gender: "F"},
	{Name: "Elena", Gender: "F"},
	{Name: "Francesca", Gender: "F"},
	{Name: "Valentina", Gender: "F"},
	{Name: "Silvia", Gender: "F"},
	{Name: "Federica", Gender: "F"},
	{Name: "Elisa", Gender: "F"},
	{Name: "Beatrice", Gender: "F"},
}

var italianSurnames = []string{
	"Rossi", "Russo", "Ferrari", "Esposito", "Bianchi",
	"Romano", "Colombo", "Ricci", "Marino", "Greco",
	"Bruno", "Gallo", "Conti", "De Luca", "Mancini",
	"Costa", "Giordano", "Rizzo", "Lombardi", "Moretti",
}

var italianPlaces = []PlaceCode{
	{Name: "Lecco", Code: "E507"},
	{Name: "Roma", Code: "H501"},
	{Name: "Milano", Code: "F205"},
	{Name: "Napoli", Code: "F839"},
	{Name: "Torino", Code: "L219"},
	{Name: "Bologna", Code: "A944"},
	{Name: "Firenze", Code: "D612"},
	{Name: "Venezia", Code: "L736"},
	{Name: "Palermo", Code: "G273"},
	{Name: "Bari", Code: "A662"},
	{Name: "Genova", Code: "D969"},
	{Name: "Verona", Code: "L781"},
	{Name: "Catania", Code: "C351"},
	{Name: "Padova", Code: "G224"},
	{Name: "Trieste", Code: "L424"},
	{Name: "Brescia", Code: "B157"},
	{Name: "Parma", Code: "G337"},
	{Name: "Modena", Code: "F257"},
	{Name: "Perugia", Code: "G547"},
	{Name: "Cagliari", Code: "B354"},
}

func generateItalianPerson(randomFunc func(int) int) PersonResult {
	nameInfo := italianNames[randomFunc(len(italianNames))]
	surname := italianSurnames[randomFunc(len(italianSurnames))]
	placeInfo := italianPlaces[randomFunc(len(italianPlaces))]

	year := 1950 + randomFunc(56)
	month := 1 + randomFunc(12)
	day := 1 + randomFunc(28)

	dateStr := fmt.Sprintf("%04d-%02d-%02d", year, month, day)

	fc := computeItalianFiscalCode(nameInfo.Name, surname, nameInfo.Gender, day, month, year, placeInfo.Code)

	return PersonResult{
		FirstName:    nameInfo.Name,
		LastName:     surname,
		Gender:       nameInfo.Gender,
		DateOfBirth:  dateStr,
		PlaceOfBirth: placeInfo.Name,
		FiscalCode:   fc,
	}
}

func computeItalianFiscalCode(firstName, lastName, gender string, day, month, year int, placeCode string) string {
	surnameCode := encodeSurname(lastName)
	nameCode := encodeFirstName(firstName)

	yearCode := fmt.Sprintf("%02d", year%100)

	monthCodes := []string{"", "A", "B", "C", "D", "E", "H", "L", "M", "P", "R", "S", "T"}
	monthCode := monthCodes[month]

	dayVal := day
	if strings.ToUpper(gender) == "F" || strings.ToUpper(gender) == "FEMALE" {
		dayVal += 40
	}
	dayCode := fmt.Sprintf("%02d", dayVal)

	partial := surnameCode + nameCode + yearCode + monthCode + dayCode + strings.ToUpper(placeCode)

	cin := computeCIN(partial)

	return partial + string(cin)
}

func encodeSurname(s string) string {
	consonants, vowels := extractLetters(s)
	combined := consonants + vowels + "XXX"
	return combined[:3]
}

func encodeFirstName(s string) string {
	consonants, vowels := extractLetters(s)
	if len(consonants) >= 4 {
		return string([]byte{consonants[0], consonants[2], consonants[3]})
	}
	combined := consonants + vowels + "XXX"
	return combined[:3]
}

func extractLetters(s string) (string, string) {
	var consonants, vowels strings.Builder
	for _, r := range strings.ToUpper(s) {
		if unicode.IsLetter(r) {
			if isVowel(r) {
				vowels.WriteRune(r)
			} else {
				consonants.WriteRune(r)
			}
		}
	}
	return consonants.String(), vowels.String()
}

func isVowel(r rune) bool {
	switch r {
	case 'A', 'E', 'I', 'O', 'U':
		return true
	default:
		return false
	}
}

func computeCIN(partial15 string) rune {
	oddMap := map[byte]int{
		'0': 1, '1': 0, '2': 5, '3': 7, '4': 9, '5': 13, '6': 15, '7': 17, '8': 19, '9': 21,
		'A': 1, 'B': 0, 'C': 5, 'D': 7, 'E': 9, 'F': 13, 'G': 15, 'H': 17, 'I': 19, 'J': 21,
		'K': 2, 'L': 4, 'M': 18, 'N': 20, 'O': 11, 'P': 3, 'Q': 6, 'R': 8, 'S': 12, 'T': 14,
		'U': 16, 'V': 10, 'W': 22, 'X': 25, 'Y': 24, 'Z': 23,
	}

	evenMap := map[byte]int{
		'0': 0, '1': 1, '2': 2, '3': 3, '4': 4, '5': 5, '6': 6, '7': 7, '8': 8, '9': 9,
		'A': 0, 'B': 1, 'C': 2, 'D': 3, 'E': 4, 'F': 5, 'G': 6, 'H': 7, 'I': 8, 'J': 9,
		'K': 10, 'L': 11, 'M': 12, 'N': 13, 'O': 14, 'P': 15, 'Q': 16, 'R': 17, 'S': 18, 'T': 19,
		'U': 20, 'V': 21, 'W': 22, 'X': 23, 'Y': 24, 'Z': 25,
	}

	sum := 0
	for i := 0; i < len(partial15); i++ {
		ch := partial15[i]
		if (i+1)%2 != 0 {
			sum += oddMap[ch]
		} else {
			sum += evenMap[ch]
		}
	}

	remainder := sum % 26
	return rune('A' + remainder)
}

// PersonCliSubcommand -----------------------------------------------------------------

type PersonCliSubcommand struct{}

func (s PersonCliSubcommand) Name() string {
	return "person"
}

func (s PersonCliSubcommand) Description() string {
	return "Generates fake person data with Codice Fiscale"
}

func (s PersonCliSubcommand) createFlagSet(outputWriter io.Writer) (*flag.FlagSet, *string) {
	personCmd := flag.NewFlagSet(s.Name(), flag.ContinueOnError)
	personCmd.SetOutput(outputWriter)
	defaultCountry := "IT"
	inputCountry := personCmd.String("country", defaultCountry, "Person country code (e.g. IT) (Only IT is supported at this time)")
	personCmd.Usage = func() {
		_, _ = fmt.Fprintf(personCmd.Output(), "Usage: generate person [-country COUNTRY_CODE]\n\nOptions:\n")
		personCmd.PrintDefaults()
	}
	return personCmd, inputCountry
}

func (s PersonCliSubcommand) Parse(args []string, outputWriter io.Writer) (Command, map[string]string, error) {
	personCmd, inputCountry := s.createFlagSet(outputWriter)
	err := personCmd.Parse(args)
	if err != nil {
		return nil, nil, err
	}
	otherArgs := map[string]string{
		"country": strings.ToUpper(*inputCountry),
	}
	return PersonCommand{}, otherArgs, nil
}

func (s PersonCliSubcommand) PrintDefaults(outputWriter io.Writer) {
	personCmd, _ := s.createFlagSet(outputWriter)
	personCmd.PrintDefaults()
}
