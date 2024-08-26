package main

import (
	"encoding/json"
	"fmt"
	"image"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/google/uuid"
)

var plainsWeights = map[string]float64{
	"cao":        0.37,
	"proteapot":  0.06,
	"berub":      0.06,
	"hiveila":    0.12,
	"sucura":     0.15,
	"sniffen":    0.05,
	"canaggi":    0.07,
	"bloobattle": 0.04,
	"oolariaf":   0.08,
}
var plainsWeightsOG = map[string]float64{
	"cao":       0.6,
	"proteapot": 0.05,
	"berub":     0.10,
	"hiveila":   0.10,
	"sucura":    0.15,
}

var mountainWeights = map[string]float64{
	"lamara":    0.6,
	"bambolt":   0.05,
	"gargrowl":  0.15,
	"celestear": 0.20,
	//"elemental": 0.15,
}

var urbanWeights = map[string]float64{
	"cao":        0.6,
	"sniffen":    0.05,
	"canaggi":    0.10,
	"bloobattle": 0.10,
	"oolariaf":   0.15,
}

var marshWeights = map[string]float64{
	"Bug":       0.5,
	"coilizard": 0.15,
	"Fairy":     0.15,
	"Whale":     0.2,
}

func genRefMap(jsonpath string) map[string]JSONwrapper {
	refMap := make(map[string]JSONwrapper)

	jsonData, err := os.ReadFile(jsonpath)
	if err != nil {
		fmt.Println("Error:", err)
	}

	json.Unmarshal(jsonData, &refMap)
	return refMap
}
func genOdds(refMap map[string]JSONwrapper) []float64 {
	c := 0.1
	var generationOdds []float64

	for key, _ := range refMap { // kinda like a mini key value cache if you think about it
		generationOdds = append(generationOdds, c+1)
		fmt.Print(key)

		c += 0.3
	} // move this out of func so doesnt called asll time
	return generationOdds
}

var typeNum = map[string]int{"Phage": 1, "Temp": 2, "Mech": 3, "Animal": 4, "Fungus": 5, "Plant": 6, "Bug": 7, "Pixie": 8, "Death": 9, "Magic": 10, "Cosmic": 11, "Mythic": 12}

// new morphs: deathclaw raptor, skeleton dinorhino, pinecone wormfairy, meteortle, torteor, asturtle

const PATTERNODDS = 1 / 2048

func genMorph(morphMap map[string]JSONwrapper, weights map[string]float64, level int) Kaizomorph {

	var chosen string

	chosen = weightedRandomSample(weights)

	pattern := false
	randpattern := rand.Float64()
	if randpattern <= PATTERNODDS {
		pattern = true
	}

	//chosen = "draginfect"
	imgpath := fmt.Sprintf("client/data/Environnement_AssetPack/sprites/%s.png", chosen) // we dont have all images yet so use temp
	//fmt.Sprintf("finalmorphs/%s.png", chosen)

	s := morphMap[chosen]
	m := &Kaizomorph{}
	mergeStructs(m, s)

	m.Name = chosen // not ideal
	m.NickName = "RandomNickname"
	m.Luck = 1 // influences random generation
	//m.Height = genHeight(*m, imgpath) // m
	m.Weight = genWeight(*m, imgpath) // kg

	m.Cuteness = genCuteness(*m) // once we have this we can normalise the Weight
	m.Weight = ((m.Weight-520.0)/(2618.0-520.0))*(1000-0) + 1.0
	m.Pattern = pattern
	m.SpeID = genSpeID()
	m.Spe = m.Spe + m.SpeID
	m.Nature = genNature(*m)
	//m.Prestige = level / 100
	if level != 100 {
		m.Level = level % 100
	} else {
		m.Level = 100
	}
	m.SizeX = 32
	m.SizeZ = 32
	fmt.Println("WEIGHT: ", m.Weight)
	m.BeamTime = 0.01 + (1 / float64(m.Weight))
	// 1-2^m.weight
	m.UUID = uuid.New().String()

	setStats(m)
	//genID
	//genHTML
	//genColor

	return *m

}

func levelUp(m Kaizomorph, n int) {

	genStats(&m, 29000)
}

// This is to make json unmarshalling easier. We then merge this inner struct with Kaizomorph
type JSONwrapper struct {
	HP       int      `json:"Hp"`
	ATK      int      `json:"Atk"`
	PDEF     int      `json:"Pdef"`
	SPA      int      `json:"Spa"`
	SPDEF    int      `json:"SPdef"`
	SPE      float32  `json:"Spe"`
	LEVEL    int      `json:"Level"`
	ITEM     string   `json:"item"`
	ABILITY  string   `json:"ability"`
	NAME     string   `json:"Name"`
	TYPES    []string `json:"type"`
	NICKNAME string   `json:"nickName"`
}

func mergeStructs(m *Kaizomorph, s JSONwrapper) {
	m.Hp = s.HP
	m.Atk = s.ATK
	m.Pdef = s.PDEF
	m.Spa = s.SPA
	m.Spdef = s.SPDEF
	m.Spe = s.SPE
	m.Level = s.LEVEL
	m.Name = s.NAME
	m.Types = s.TYPES
	m.NickName = s.NICKNAME

}

// ReferencImplementation
// should moves be enums?
// put related field types together for efficient storage (bool should be next to bool, for example)
type MorphColor struct {
	r int
	g int
	b int
	a int
}
type Kaizomorph struct {
	Image []byte
	Hp    int
	Atk   int
	Pdef  int
	Spa   int
	Spdef int
	Spe   float32
	SpeID float32
	Luck  int
	Level int // some default low value like 1 or 5 in json file
	//Item  //GameItem // set to null or empty if no item
	Types  []string
	Nature string
	//Moves  []Move
	Stam   []int32
	Status string
	Exp    int
	// generate Height and Weight based on convex hull of image
	Height  float32
	Weight  float32
	MorphID int
	OwnerID int
	// geolocation where user generated
	FoundAt  string
	DateMet  int
	Name     string
	NickName string
	Species  string
	// effort values effect the final value of a stat
	Hp_evs       int
	Atk_evs      int
	Pdef_evs     int
	Spa_evs      int
	SPdef_evs    int
	Spe_evs      int
	XPyield      int
	EVyield      int
	Kills        int
	Deaths       int
	UnusedInt    int
	UnusedString string
	UnusedBool   bool
	UnusedArray  []int

	Group    string
	Color    MorphColor
	Pattern  bool
	Cuteness int
	Flags    []bool // pass into functions for validation
	/* EVERYTHING BELOW THIS LINE IS SPECIFIC TO THIS VERSION OF THE GAME */
	Pos      *Vector3
	BeamTime float64
	SizeX    float64
	SizeZ    float64
	UUID     string
}

func setStats(m *Kaizomorph) {
	m.Hp = (((2*m.Hp + (m.Hp_evs / 4)) * m.Level) / 100) + m.Level + 10
	m.Atk = (((((2*m.Atk + (m.Atk_evs / 4)) * m.Level) / 100) + 5) * 1)
	m.Pdef = (((((2*m.Pdef + (m.Pdef_evs / 4)) * m.Level) / 100) + 5) * 1)
	m.Spa = (((((2*m.Spa + (m.Spa_evs / 4)) * m.Level) / 100) + 5) * 1)
	m.Spdef = (((((2*m.Spdef + (m.SPdef_evs / 4)) * m.Level) / 100) + 5) * 1)
	m.Spe = float32((((((2*int(m.Spe) + (m.Spe_evs / 4)) * m.Level) / 100) + 5) * 1)) // + m.SpeID
}

func genStats(m *Kaizomorph, EV float64) {
	m.Hp = int(2.0*float64(m.Hp) + EV/4.0*float64(m.Level)/100.0 + float64(m.Level)) // +nature% +natureID // personality?
	m.Atk = int(2.0*float64(m.Atk) + EV/4.0*float64(m.Level)/100.0 + float64(m.Level))
	m.Pdef = int(2.0*float64(m.Pdef) + EV/4.0*float64(m.Level)/100.0 + float64(m.Level))
	m.Spa = int(2.0*float64(m.Spa) + EV/4.0*float64(m.Level)/100.0 + float64(m.Level))
	m.Spdef = int(2.0*float64(m.Spdef) + EV/4.0*float64(m.Level)/100.0 + float64(m.Level))
	m.Spe = float32(2.0*float64(m.Spe) + EV/4.0*float64(m.Level)/100.0 + float64(m.Level))
}

func genHTML(m *Kaizomorph) string {

	rawHTML := `<html>
				Name: %s 
				Age: %d
				Height: %.2f
				Is student? %t`

	cleanHTML := fmt.Sprintf(rawHTML, m.Name, m.Level)
	return cleanHTML
}

func genSpeID() float32 {

	// Decimals appended to Speed to reduce Speed ties and give players a reason to collect
	// The chance of any two 3 decimal ID's being the same is 1/1000, adding that by
	// whatever chance there is to have the same base Speed.
	// (it is not 1/496 because Speeds are not uniformly distributed, but something like that)
	randomFloat := rand.Float32()
	truncated := int(randomFloat * 100)
	SpeID := float32(truncated) / 100.0
	return SpeID
}

func genCuteness(m Kaizomorph) int {

	// get the range by calling genWeight on smallest and largest
	// smallest = 520 (paraseed)
	// largest = 2618 (pomegladon)
	// NOTE that we normalise the number of pixels in an image, not the actual Weight
	// the Weight itself has to be normalised so normalising twice can create bad distribution of values

	// Should we factor in height so that particularly tall skinny morphs arent considered too cute?

	W := float64(m.Weight)
	W_min := 720.0
	W_max := 2618.0
	norm := 5.0 / (((W-W_min)/(W_max-W_min))*(5.0-1.0) + 1.0) // Natures between 1-5

	intPart := math.Floor(norm)
	fracPart := norm - intPart
	if fracPart > 0.5 {
		Cuteness := intPart + 1.0
		if Cuteness > 5.0 {
			Cuteness = 5.0
			return int(Cuteness)
		}

		return int(Cuteness)
	}
	Cuteness := int(intPart)
	if Cuteness > 5 {
		Cuteness = 5
		return Cuteness
	}

	return Cuteness
}

func genNature(m Kaizomorph) string {
	// user ranking on Cuteness?
	switch m.Cuteness { //represented by stars
	case 1:
		Natures := [5]string{"Hardy", "Hasty", "Lonely", "Lax", "Rash"}
		randomIndex := rand.Intn(5)
		Nature := Natures[randomIndex]
		return Nature

	case 2:
		Natures := [5]string{"Brave", "Naive", "Quiet", "Careful", "Serious"}
		randomIndex := rand.Intn(5)
		Nature := Natures[randomIndex]
		return Nature

	case 3:
		Natures := [5]string{"Mild", "Calm", "Bold", "Docile", "Adamant"}
		randomIndex := rand.Intn(5)
		Nature := Natures[randomIndex]
		return Nature

	case 4:
		Natures := [5]string{"Bashful", "Jolly", "Impish", "Gentle", "Relaxed"}
		randomIndex := rand.Intn(5)
		Nature := Natures[randomIndex]
		return Nature

	case 5:
		Natures := [5]string{"Naughty", "Timid", "Sassy", "Quirky", "Modest"}
		randomIndex := rand.Intn(5)
		Nature := Natures[randomIndex]
		return Nature

	}
	return "Serious"

}

// todo
func genMoves(m Kaizomorph) []string {
	moves := []string{}
	moves = append(moves, "Sen's Paradox")
	return moves
}

func genLocation(m *Kaizomorph, ip string) {
	m.FoundAt = ip
	m.DateMet = time.Now().Day()
}

func genImage(m *Kaizomorph, pattern bool) {
	displayImage := "/some.png" // use image for checksum,
	fmt.Print(displayImage)

	if pattern {
		m.Image = []byte{1} //"patternpath" or genPattern()
		return
	}
	//m.Image = ""

	imageByteData, err := ioutil.ReadFile(displayImage)
	if err != nil {
		fmt.Println("Error decoding the image:", err)
	}
	m.Image = imageByteData
}

func genWeight(m Kaizomorph, path string) float32 {

	inputFile, err := os.Open(path)
	if err != nil {
		return -1
	}
	defer inputFile.Close()

	img, _, err := image.Decode(inputFile)
	if err != nil {
		return -1
	}

	numPixels := 0

	bounds := img.Bounds()

	// loop forward
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			_, _, _, a := img.At(x, y).RGBA()

			if a == 0 {
				continue

			}
			numPixels += 1

		}
	}

	Weight := numPixels
	fmt.Println("NUMPIXELS:", Weight)
	return float32(Weight)
}

func genHeight(m Kaizomorph, path string) float32 {

	inputFile, err := os.Open(path)
	if err != nil {
		fmt.Print(err)
	}
	defer inputFile.Close()

	img, _, err := image.Decode(inputFile)
	if err != nil {
		fmt.Print(err)
	}

	var headPixel Pixel
	var headPixelLoc Point
	var tailPixel Pixel
	var tailPixelLoc Point

	bounds := img.Bounds()

	// loop forward
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()

			if a == 0 {
				continue
			}
			if headPixelLoc == (Point{}) {
				headPixelLoc.x = x
				headPixelLoc.y = y
				tailPixel.r = r
				tailPixel.g = g
				tailPixel.b = b
				tailPixel.a = a

			}

		}
	}

	// loop backward

	for y := bounds.Max.Y; y > bounds.Min.Y; y-- {
		for x := bounds.Max.X; x > bounds.Min.X; x-- {
			r, g, b, a := img.At(x, y).RGBA()

			if a == 1 {
				continue
			}
			if tailPixelLoc == (Point{}) {
				tailPixelLoc.x = x
				tailPixelLoc.y = y
				headPixel.r = r
				headPixel.g = g
				headPixel.b = b
				headPixel.a = a

			}

		}
	}

	Height := tailPixelLoc.y - headPixelLoc.y

	randomFloat := -0.5 + rand.Float64()*(1.5)

	// so we have a pixel Height but it needs to look like a real Height
	// we'll say 1 pixel is
	// pomeg - 63 therefore is 6.3 meters
	// smallest - 30-49 or so, about 0.4 meters
	// let's say the biggest you can be is 6.8 meters
	// every pixel adds.. 30/6.4 = .45 meters
	min := 0.4
	max := 6.9
	normHeight := ((float64(Height)-30.0)/(63.0-30.0))*(max-min) + min

	return float32(normHeight) + float32(randomFloat)
}
