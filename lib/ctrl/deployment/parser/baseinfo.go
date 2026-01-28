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

package parser

import (
	"strconv"

	"github.com/beevik/etree"
)

func (this *Parser) getOrder(element *etree.Element) int64 {
	if element == nil {
		return 0
	}
	orderAttr := element.SelectAttr("order")
	if orderAttr == nil {
		return 0
	}
	order, err := strconv.ParseInt(orderAttr.Value, 10, 64)
	if err != nil {
		this.conf.GetLogger().Debug("unable to parse element senergy:order as int64", "error", err, "element", *element)
		return 0
	}
	return order
}
