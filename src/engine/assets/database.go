/******************************************************************************/
/* database.go                                                                */
/******************************************************************************/
/*                           This file is part of:                            */
/*                                KAIJU ENGINE                                */
/*                          https://kaijuengine.org                           */
/******************************************************************************/
/* MIT License                                                                */
/*                                                                            */
/* Copyright (c) 2023-present Kaiju Engine authors (AUTHORS.md).              */
/* Copyright (c) 2015-present Brent Farris.                                   */
/*                                                                            */
/* May all those that this source may reach be blessed by the LORD and find   */
/* peace and joy in life.                                                     */
/* Everyone who drinks of this water will be thirsty again; but whoever       */
/* drinks of the water that I will give him shall never thirst; John 4:13-14  */
/*                                                                            */
/* Permission is hereby granted, free of charge, to any person obtaining a    */
/* copy of this software and associated documentation files (the "Software"), */
/* to deal in the Software without restriction, including without limitation  */
/* the rights to use, copy, modify, merge, publish, distribute, sublicense,   */
/* and/or sell copies of the Software, and to permit persons to whom the      */
/* Software is furnished to do so, subject to the following conditions:       */
/*                                                                            */
/* The above copyright, blessing, biblical verse, notice and                  */
/* this permission notice shall be included in all copies or                  */
/* substantial portions of the Software.                                      */
/*                                                                            */
/* THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS    */
/* OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF                 */
/* MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.     */
/* IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY       */
/* CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT  */
/* OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE      */
/* OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.                              */
/******************************************************************************/

package assets

import (
	"kaiju/platform/filesystem"
	"kaiju/platform/profiler/tracing"
)

type Database struct {
	EditorContext EditorContext
	cache         map[string][]byte
}

func NewDatabase() Database {
	return Database{
		cache: make(map[string][]byte),
	}
}

func (a *Database) ToFilePath(key string) string { return a.toContentPath(key) }

func (a *Database) Cache(key string, data []byte) { a.cache[key] = data }
func (a *Database) CacheRemove(key string)        { delete(a.cache, key) }
func (a *Database) CacheClear()                   { clear(a.cache) }

func (a *Database) ReadText(key string) (string, error) {
	defer tracing.NewRegion("AssetDatabase.ReadText: " + key).End()
	if data, ok := a.cache[key]; ok {
		return string(data), nil
	}
	return filesystem.ReadTextFile(a.toContentPath(key))
}

func (a *Database) Read(key string) ([]byte, error) {
	defer tracing.NewRegion("AssetDatabase.Read: " + key).End()
	if data, ok := a.cache[key]; ok {
		return data, nil
	}
	return filesystem.ReadFile(a.toContentPath(key))
}

func (a *Database) Exists(key string) bool {
	if _, ok := a.cache[key]; ok {
		return true
	}
	return filesystem.FileExists(a.toContentPath(key))
}

func (a *Database) Destroy() {

}
