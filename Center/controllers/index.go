// Copyright 2013 Beego Samples authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package controllers

type IndexController struct {
	BaseController
}

//默认首页
func (this *IndexController) Get() {
	this.Display("index")
}

//ws连接
func (this *IndexController) Join() {

	uname := this.GetString("uname")
	tech := this.GetString("tech")

	// Check valid.
	if len(uname) == 0 {
		this.Redirect("/", 302)
		return
	}

	switch tech {
	case "longpolling":
		this.Redirect("/lp?uname="+uname, 302)
	case "websocket":
		this.Redirect("/ws?uname="+uname, 302)
	default:
		this.Redirect("/", 302)
	}

	return
}
