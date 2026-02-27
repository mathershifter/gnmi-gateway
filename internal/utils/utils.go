// Copyright 2020 Netflix Inc
// Author: Colin McIntosh (colin@netflix.com)
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

package utils

import (
	"fmt"
	"sort"
	"strings"

	"github.com/openconfig/gnmi/proto/gnmi"
)

// GNMINotificationPrettyString returns a human-readable string for
// a gnmi.Notification.
func GNMINotificationPrettyString(notification *gnmi.Notification) string {
	if notification == nil {
		return "<nil>"
	}

	var pretty strings.Builder
	pretty.WriteString("{ ")
	fmt.Fprintf(&pretty, "timestamp:%d ", notification.Timestamp)
	fmt.Fprintf(&pretty, "prefix:%s ", PathToXPath(notification.Prefix))
	// if notification.Alias != "" {
	// 	pretty += fmt.Sprintf("alias:%v ", notification.Alias)
	// }
	if len(notification.Update) > 0 {
		pretty.WriteString("update:[ ")
		for _, u := range notification.Update {
			pretty.WriteString(GNMIUpdatePrettyString(u) + " ")
		}
		pretty.WriteString("] ")
	}
	if len(notification.Delete) > 0 {
		pretty.WriteString("delete:[ ")
		for _, d := range notification.Delete {
			pretty.WriteString(PathToXPath(d) + " ")
		}
		pretty.WriteString("] ")
	}
	if notification.Atomic {
		pretty.WriteString("atomic:true ")
	}

	pretty.WriteString("}")
	return pretty.String()
}

// GNMIUpdatePrettyString returns a human-readable string for a gnmi.Update.
func GNMIUpdatePrettyString(update *gnmi.Update) string {
	if update == nil {
		return "<nil>"
	}

	pretty := "{ "
	pretty += fmt.Sprintf("path:%s ", PathToXPath(update.Path))
	pretty += fmt.Sprintf("val:%v ", update.Val.GetValue())
	if update.Duplicates != 0 {
		pretty += fmt.Sprintf("duplicates:%d ", update.Duplicates)
	}
	pretty += "}"
	return pretty
}

// PathToXPath returns an XPath-style string for a gnmi.Path.
func PathToXPath(path *gnmi.Path) string {
	var xPath strings.Builder

	if path.Origin != "" {
		xPath.WriteString("/" + path.Origin)
	}

	if path.Target != "" {
		xPath.WriteString("/" + path.Target)
	}

	for _, elem := range path.Elem {
		if elem.Name != "" {
			xPath.WriteString("/" + elem.Name)
		}

		if len(elem.Key) > 0 {
			keys := make([]string, 0, len(elem.Key))
			for k := range elem.Key {
				keys = append(keys, k)
			}
			sort.Strings(keys)

			pairs := make([]string, 0, len(elem.Key))
			for _, key := range keys {
				pairs = append(pairs, key+"="+elem.Key[key])
			}
			xPath.WriteString("[" + strings.Join(pairs, " ") + "]")
		}
	}
	return xPath.String()
}

func GetNumberValues(tv *gnmi.TypedValue) (float64, bool) {
	if tv != nil && tv.Value != nil {
		switch tv.Value.(type) {
		case *gnmi.TypedValue_StringVal:
			return 0, false
		case *gnmi.TypedValue_IntVal:
			return float64(tv.GetIntVal()), true
		case *gnmi.TypedValue_UintVal:
			return float64(tv.GetUintVal()), true
		case *gnmi.TypedValue_BoolVal:
			if tv.GetBoolVal() {
				return 1, true
			} else {
				return 0, false
			}
		case *gnmi.TypedValue_FloatVal:
			//nolint:staticcheck
			return float64(tv.GetFloatVal()), true
		case *gnmi.TypedValue_LeaflistVal:
			return 0, false
		case *gnmi.TypedValue_BytesVal:
			return 0, false
		default:
			return 0, false
		}
	}
	return 0, false
}
