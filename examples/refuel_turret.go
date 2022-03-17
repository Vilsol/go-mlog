package main

import "github.com/Vilsol/go-mlog/m"

//Requires a unit to be flagged
func main() {
	flag := 2
	unitType := m.UMono
	const turretName = "duo1"
	ammoType := m.ITCopper
	m.UnitBind(unitType)
	if m.CurUnit.GetFlag() != flag {
		return
	}
	item := m.CurUnit.GetFirstItem()
	if item != ammoType { //go to core and pick up
		coreX, coreY, _, core := m.UnitLocateBuilding(m.BCore, false)
		m.UnitApproach(coreX, coreY, m.CurUnit.GetRange()-1)
		unitCapacity := m.CurUnit.GetTotalItems()
		m.UnitItemDrop(core, unitCapacity)           //drop non-ammo item
		m.UnitItemTake(core, ammoType, unitCapacity) //take ammo item
	} else { // go to turret and refuel
		turret := m.GetTurret(turretName)
		m.UnitApproach(turret.GetX(), turret.GetY(), m.CurUnit.GetRange()-1)
		m.UnitItemDrop(turret, m.CurUnit.GetTotalItems())
	}
}
