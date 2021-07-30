/**
 * Copyright (C) 2021 Yi Fan Song <yfsong00@gmail.com>
 *
 * This file is part of Goani.
 *
 * Goani is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Goani is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with Goani.  If not, see <https://www.gnu.org/licenses/>.
 **/

package main

import "os"

func createFile(filename string) (*os.File, error) {
	f, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	return f, nil
}
