package m

func Read(memory string, position int) int {
	return 0
}

func Write(value int, memory string, position int) {
}

func DrawClear(r int, g int, b int) {
}

func DrawColor(r int, g int, b int, a int) {
}

func DrawStroke(width int) {
}

func DrawLine(x1 int, y1 int, x2 int, y2 int) {
}

func DrawRect(x int, y int, width int, height int) {
}

func DrawLineRect(x int, y int, width int, height int) {
}

func DrawPoly(x int, y int, sides int, radius float32, rotation float32) {
}

func DrawLinePoly(x int, y int, sides int, radius float32, rotation float32) {
}

func DrawTriangle(x1 int, y1 int, x2 int, y2 int, x3 int, y3 int) {
}

func DrawImage(x int, y int, image string, size float32, rotation float32) {
}

func PrintFlush(targetMessage string) {
}

func DrawFlush(targetDisplay string) {
}

func GetLink(output string, address string) {
}

func ControlEnabled(target string, enabled int) {
}

func ControlShoot(target string, x int, y int, shoot int) {
}

func ControlShootP(target string, unit int, shoot int) {
}

func ControlConfigure(target string, configuration int) {
}

type RadarTarget = int

const (
	RTAny = RadarTarget(iota)
	RTEnemy
	RTAlly
	RTPlayer
	RTAttacker
	RTFlying
	RTBoss
	RTGround
)

type RadarSort = int

const (
	RSDistance = RadarSort(iota)
	RSHealth
	RSShield
	RSArmor
	RSMaxHealth
)

func Radar(from string, target1 RadarTarget, target2 RadarTarget, target3 RadarTarget, sortOrder int, sort RadarSort, output string) {
}

func Sensor(variable string, sensor string, from string) {
}

// TODO Operations

func End() {
}

func Floor(number float64) int {
	return 0
}

func Random(max float64) float64 {
	return 0
}

func Sleep(millis int) {
}
