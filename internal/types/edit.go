/*
 * © 2022-2025 Snyk Limited
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package types

type TextEdit struct {

	/**
	 * The range of the text document to be manipulated. To insert
	 * text into a document create a range where start === end.
	 */
	Range Range

	/**
	 * The string to be inserted. For delete operations use an
	 * empty string.
	 */
	NewText string
}

type WorkspaceEdit struct {
	/**
	 * Holds changes to existing resources, keyed on the affected file path.
	 */
	Changes map[string][]TextEdit
}
