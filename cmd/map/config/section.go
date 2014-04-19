// Copyright 2009  The "config" Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

// import "fmt"

// AddSection adds a new section to the configuration.
//
// If the section is nil then uses the section by default which it's already
// created.
func (c *Config) AddSection(section string) {
	// DEFAULT_SECTION
	if section == "" {
		return
	}

	// a array
	if secs, ok := c.sections[section]; ok {
		secs = append(secs, make(map[string]*tValue))
		c.sections[section] = secs
		// fmt.Println("this section is a array:", section)
	} else {
		// already in
		if sec, ok := c.section[section]; ok {
			c.section[section] = nil
			secs := make([]Section, 0, 4)
			secs = append(secs, sec)
			secs = append(secs, make(map[string]*tValue))
			c.sections[section] = secs
			// fmt.Println("already hava the section, change to array:", section, len(c.sections[section]))
		} else {
			c.section[section] = make(map[string]*tValue)
			// fmt.Println("add to section:", section)
		}
	}
}

// RemoveSection removes a section from the configuration.
// It returns true if the section was removed, and false if section did not exist.
// func (c *Config) RemoveSection(section string) bool {
// 	_, ok := c.data[section]

// 	// Default section cannot be removed.
// 	if !ok || section == DEFAULT_SECTION {
// 		return false
// 	}

// 	for o, _ := range c.data[section] {
// 		delete(c.data[section], o) // *value
// 	}
// 	delete(c.data, section)

// 	delete(c.lastIdOption, section)
// 	delete(c.idSection, section)

// 	return true
// }

// HasSection checks if the configuration has the given section.
// (The default section always exists.)
func (c *Config) HasSection(section string) bool {
	_, secs := c.sections[section]
	_, sec := c.section[section]
	return secs || sec
}

// Sections returns the list of sections in the configuration.
// (The default section always exists).
// func (c *Config) Sections() (sections []string) {
// 	sections = make([]string, len(c.idSection))
// 	pos := 0 // Position in sections

// 	for i := 0; i < c.lastIdSection; i++ {
// 		for section, id := range c.idSection {
// 			if id == i {
// 				sections[pos] = section
// 				pos++
// 			}
// 		}
// 	}

// 	return sections
// }
