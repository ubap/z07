	go-tibia - A Man-in-the-Middle (MITM) proxy for the Tibia MMO
	Copyright (C) 2025 Jakub Trzebiatowski

	This program is free software: you can redistribute it and/or modify
	it under the terms of the GNU General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	This program is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU General Public License for more details.

	You should have received a copy of the GNU General Public License
	along with this program.  If not, see <https://www.gnu.org/licenses/>.


A Man-in-the-Middle (MITM) proxy for the Tibia MMO, built to explore network programming and protocol reverse engineering in Go.

----
Package structure (outdated)

    Login --> Packets
    Login --> Protocol

    Packets --> Protocol

    Protocol --> Nothing
    Model --> Nothing