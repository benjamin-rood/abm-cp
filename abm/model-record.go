package abm

func (m *Model) cpPreyRecordCopy() map[string]ColourPolymorphicPrey {
	defer m.rcpPreyRW.RUnlock()
	m.rcpPreyRW.RLock()
	var record = make(map[string]ColourPolymorphicPrey)
	for k, v := range m.recordCPP {
		record[k] = v
	}
	return record
}

func (m *Model) cpPreyRecordAssignValue(key string, value ColourPolymorphicPrey) error {
	defer m.rcpPreyRW.Unlock()
	m.rcpPreyRW.Lock()
	m.recordCPP[key] = value
	return nil
}

func (m *Model) vpRecordCopy() map[string]VisualPredator {
	defer m.rvpRW.RUnlock()
	m.rvpRW.RLock()
	var record = make(map[string]VisualPredator)
	for k, v := range m.recordVP {
		record[k] = v
	}
	return record
}

func (m *Model) vpRecordAssignValue(key string, value VisualPredator) error {
	defer m.rvpRW.Unlock()
	m.rvpRW.Lock()
	m.recordVP[key] = value
	return nil
}
