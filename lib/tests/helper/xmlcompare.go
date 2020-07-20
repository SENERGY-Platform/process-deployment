/*
 * Copyright 2020 InfAI (CC SES)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package helper

import (
	"encoding/xml"
	"regexp"
)

func XmlIsEqual(a string, b string) (equal bool, err error) {
	aNormal, err := normalizeXml(a)
	if err != nil {
		return false, err
	}
	bNormal, err := normalizeXml(b)
	if err != nil {
		return false, err
	}
	return aNormal == bNormal, nil
}

func normalizeXml(xmlStr string) (result string, err error) {
	var normalized Node
	err = xml.Unmarshal([]byte(xmlStr), &normalized)
	if err != nil {
		return "", err
	}
	recursiveReplace(&normalized)
	buf, err := xml.Marshal(normalized)
	return string(buf), err
}

type Node struct {
	XMLName  xml.Name
	Attrs    []xml.Attr `xml:",any,attr"`
	Text     string     `xml:",chardata"`
	Children []Node     `xml:",any"`
}

func recursiveReplace(n *Node) {
	n.Text = regexp.MustCompile(`\n\s*`).ReplaceAllString(n.Text, "")
	for i := range n.Children {
		recursiveReplace(&n.Children[i])
	}
}
