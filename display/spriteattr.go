/*H**********************************************************************
* FILENAME :        spriteattr.go
*
* PACKAGE :			display
*
* AUTHOR :    Moshe Nahmias       LAST CHANGE :    04 Jan 2017
*
*H*/

package display

// SpriteAboveBackground priority
const SpriteAboveBackground byte = 0

// SpriteBelowBackground priority
const SpriteBelowBackground byte = 1

// ObjPalette0 represents OBP0
const ObjPalette0 byte = 0

// ObjPalette1 represents OBP1
const ObjPalette1 byte = 1

// SpriteAttr from OAM
type SpriteAttr [4]byte

// SpriteAttrs from OAM
type SpriteAttrs []*SpriteAttr

func (s SpriteAttrs) Len() int {
	return len(s)
}

func (s SpriteAttrs) Less(i, j int) bool {
	return s[i].coordinateX() < s[j].coordinateX() ||
		(s[i].coordinateX() == s[j].coordinateX() && s[i].tileID() < s[j].tileID())
}

func (s SpriteAttrs) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s *SpriteAttr) coordinateY() byte {
	return s[0]
}

func (s *SpriteAttr) coordinateX() byte {
	return s[1]
}

func (s *SpriteAttr) tileID() byte {
	return s[2]
}

func (s *SpriteAttr) priority() byte {
	return s[3] >> 7
}

func (s *SpriteAttr) flipY() bool {
	return s[3]&0x40 == 0x40
}

func (s *SpriteAttr) flipX() bool {
	return s[3]&0x20 == 0x20
}

func (s *SpriteAttr) palette() byte {
	return (s[3] << 3) >> 7
}
