/* Copyright 2014 Ooyala, Inc. All rights reserved.
 *
 * This file is licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
 * except in compliance with the License. You may obtain a copy of the License at
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License is
 * distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and limitations under the License.
 */

package builder

import (
	. "atlantis/common"
	"io"
)

var DefaultBuilder Builder

type Builder interface {
	// (task, repo, root, sha)
	Build(*Task, string, string, string) (io.ReadCloser, error)
	// (task, repo, root, sha, user, password)
	AuthenticatedBuild(*Task, string, string, string, string, string)
}
