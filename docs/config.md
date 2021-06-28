---
title: API Documentation
linktitle: API Documentation
description: Reference of the jx-promote configuration
weight: 10
---
<p>Packages:</p>
<ul>
<li>
<a href="#updatebot.jenkins-x.io%2fv1alpha1">updatebot.jenkins-x.io/v1alpha1</a>
</li>
</ul>
<h2 id="updatebot.jenkins-x.io/v1alpha1">updatebot.jenkins-x.io/v1alpha1</h2>
<p>
<p>Package v1alpha1 is the v1alpha1 version of the API.</p>
</p>
Resource Types:
<ul><li>
<a href="#updatebot.jenkins-x.io/v1alpha1.UpdateConfig">UpdateConfig</a>
</li></ul>
<h3 id="updatebot.jenkins-x.io/v1alpha1.UpdateConfig">UpdateConfig
</h3>
<p>
<p>UpdateConfig defines the update rules</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>apiVersion</code></br>
string</td>
<td>
<code>
updatebot.jenkins-x.io/v1alpha1
</code>
</td>
</tr>
<tr>
<td>
<code>kind</code></br>
string
</td>
<td><code>UpdateConfig</code></td>
</tr>
<tr>
<td>
<code>metadata</code></br>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.13/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
<em>(Optional)</em>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code></br>
<em>
<a href="#updatebot.jenkins-x.io/v1alpha1.UpdateConfigSpec">
UpdateConfigSpec
</a>
</em>
</td>
<td>
<p>Spec holds the update rule specifications</p>
<br/>
<br/>
<table>
<tr>
<td>
<code>rules</code></br>
<em>
<a href="#updatebot.jenkins-x.io/v1alpha1.Rule">
[]Rule
</a>
</em>
</td>
<td>
<p>Rules defines the change rules</p>
</td>
</tr>
</table>
</td>
</tr>
</tbody>
</table>
<h3 id="updatebot.jenkins-x.io/v1alpha1.Change">Change
</h3>
<p>
(<em>Appears on:</em>
<a href="#updatebot.jenkins-x.io/v1alpha1.Rule">Rule</a>)
</p>
<p>
<p>Change the kind of change to make on a repository</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>command</code></br>
<em>
<a href="#updatebot.jenkins-x.io/v1alpha1.Command">
Command
</a>
</em>
</td>
<td>
<p>Command runs a shell command</p>
</td>
</tr>
<tr>
<td>
<code>go</code></br>
<em>
<a href="#updatebot.jenkins-x.io/v1alpha1.GoChange">
GoChange
</a>
</em>
</td>
<td>
<p>Go for go lang based dependency upgrades</p>
</td>
</tr>
<tr>
<td>
<code>regex</code></br>
<em>
<a href="#updatebot.jenkins-x.io/v1alpha1.Regex">
Regex
</a>
</em>
</td>
<td>
<p>Regex a regex based modification</p>
</td>
</tr>
<tr>
<td>
<code>versionStream</code></br>
<em>
<a href="#updatebot.jenkins-x.io/v1alpha1.VersionStreamChange">
VersionStreamChange
</a>
</em>
</td>
<td>
<p>VersionStream updates the charts in a version stream repository</p>
</td>
</tr>
<tr>
<td>
<code>versionTemplate</code></br>
<em>
string
</em>
</td>
<td>
<p>VersionTemplate an optional template if the version is coming from a previous Pull Request SHA</p>
</td>
</tr>
</tbody>
</table>
<h3 id="updatebot.jenkins-x.io/v1alpha1.Command">Command
</h3>
<p>
(<em>Appears on:</em>
<a href="#updatebot.jenkins-x.io/v1alpha1.Change">Change</a>)
</p>
<p>
<p>Command runs a command line program</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code></br>
<em>
string
</em>
</td>
<td>
<p>Name the name of the command</p>
</td>
</tr>
<tr>
<td>
<code>args</code></br>
<em>
[]string
</em>
</td>
<td>
<p>Args the command line arguments</p>
</td>
</tr>
<tr>
<td>
<code>env</code></br>
<em>
<a href="#updatebot.jenkins-x.io/v1alpha1.EnvVar">
[]EnvVar
</a>
</em>
</td>
<td>
<p>Env the environment variables to pass into the command</p>
</td>
</tr>
</tbody>
</table>
<h3 id="updatebot.jenkins-x.io/v1alpha1.EnvVar">EnvVar
</h3>
<p>
(<em>Appears on:</em>
<a href="#updatebot.jenkins-x.io/v1alpha1.Command">Command</a>)
</p>
<p>
<p>EnvVar the environment variable</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code></br>
<em>
string
</em>
</td>
<td>
<p>Name the name of the environment variable</p>
</td>
</tr>
<tr>
<td>
<code>value</code></br>
<em>
string
</em>
</td>
<td>
<p>Value the value of the environment variable</p>
</td>
</tr>
</tbody>
</table>
<h3 id="updatebot.jenkins-x.io/v1alpha1.GoChange">GoChange
</h3>
<p>
(<em>Appears on:</em>
<a href="#updatebot.jenkins-x.io/v1alpha1.Change">Change</a>)
</p>
<p>
<p>GoChange for upgrading go dependencies</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>owner</code></br>
<em>
[]string
</em>
</td>
<td>
<p>Owners the git owners to query</p>
</td>
</tr>
<tr>
<td>
<code>repositories</code></br>
<em>
<a href="#updatebot.jenkins-x.io/v1alpha1.Pattern">
Pattern
</a>
</em>
</td>
<td>
<p>Repositories the repositories to match</p>
</td>
</tr>
<tr>
<td>
<code>package</code></br>
<em>
string
</em>
</td>
<td>
<p>Package the text in the go.mod to filter on to perform an upgrade</p>
</td>
</tr>
<tr>
<td>
<code>upgradePackages</code></br>
<em>
<a href="#updatebot.jenkins-x.io/v1alpha1.Pattern">
Pattern
</a>
</em>
</td>
<td>
<p>UpgradePackages the packages to upgrade</p>
</td>
</tr>
<tr>
<td>
<code>noPatch</code></br>
<em>
bool
</em>
</td>
<td>
<p>NoPatch disables patch upgrades so we can import to new minor releases</p>
</td>
</tr>
</tbody>
</table>
<h3 id="updatebot.jenkins-x.io/v1alpha1.Pattern">Pattern
</h3>
<p>
(<em>Appears on:</em>
<a href="#updatebot.jenkins-x.io/v1alpha1.GoChange">GoChange</a>, 
<a href="#updatebot.jenkins-x.io/v1alpha1.VersionStreamChange">VersionStreamChange</a>)
</p>
<p>
<p>Pattern for matching strings</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code></br>
<em>
string
</em>
</td>
<td>
<p>Name</p>
</td>
</tr>
<tr>
<td>
<code>include</code></br>
<em>
[]string
</em>
</td>
<td>
<p>Includes patterns to include in changing</p>
</td>
</tr>
<tr>
<td>
<code>exclude</code></br>
<em>
[]string
</em>
</td>
<td>
<p>Excludes patterns to exclude from upgrading</p>
</td>
</tr>
</tbody>
</table>
<h3 id="updatebot.jenkins-x.io/v1alpha1.Regex">Regex
</h3>
<p>
(<em>Appears on:</em>
<a href="#updatebot.jenkins-x.io/v1alpha1.Change">Change</a>)
</p>
<p>
<p>Regex a regex based modification</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>pattern</code></br>
<em>
string
</em>
</td>
<td>
<p>Pattern the regex pattern to apply</p>
</td>
</tr>
<tr>
<td>
<code>files</code></br>
<em>
[]string
</em>
</td>
<td>
<p>Globs the files to apply this to</p>
</td>
</tr>
</tbody>
</table>
<h3 id="updatebot.jenkins-x.io/v1alpha1.Rule">Rule
</h3>
<p>
(<em>Appears on:</em>
<a href="#updatebot.jenkins-x.io/v1alpha1.UpdateConfigSpec">UpdateConfigSpec</a>)
</p>
<p>
<p>Rule specifies a set of repositories and changes</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>urls</code></br>
<em>
[]string
</em>
</td>
<td>
<p>URLs the git URLs of the repositories to create a Pull Request on</p>
</td>
</tr>
<tr>
<td>
<code>changes</code></br>
<em>
<a href="#updatebot.jenkins-x.io/v1alpha1.Change">
[]Change
</a>
</em>
</td>
<td>
<p>Changes the changes to perform on the repositories</p>
</td>
</tr>
<tr>
<td>
<code>fork</code></br>
<em>
bool
</em>
</td>
<td>
<p>Fork if we should create the pull request from a fork of the repository</p>
</td>
</tr>
</tbody>
</table>
<h3 id="updatebot.jenkins-x.io/v1alpha1.UpdateConfigSpec">UpdateConfigSpec
</h3>
<p>
(<em>Appears on:</em>
<a href="#updatebot.jenkins-x.io/v1alpha1.UpdateConfig">UpdateConfig</a>)
</p>
<p>
<p>UpdateConfigSpec defines the rules to perform when updating.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>rules</code></br>
<em>
<a href="#updatebot.jenkins-x.io/v1alpha1.Rule">
[]Rule
</a>
</em>
</td>
<td>
<p>Rules defines the change rules</p>
</td>
</tr>
</tbody>
</table>
<h3 id="updatebot.jenkins-x.io/v1alpha1.VersionStreamChange">VersionStreamChange
</h3>
<p>
(<em>Appears on:</em>
<a href="#updatebot.jenkins-x.io/v1alpha1.Change">Change</a>)
</p>
<p>
<p>VersionStreamChange for upgrading versions in a version stream</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>Pattern</code></br>
<em>
<a href="#updatebot.jenkins-x.io/v1alpha1.Pattern">
Pattern
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>kind</code></br>
<em>
string
</em>
</td>
<td>
<p>Kind the kind of resources to change (charts, git, package etc)</p>
</td>
</tr>
</tbody>
</table>
<hr/>
<p><em>
Generated with <code>gen-crd-api-reference-docs</code>
on git commit <code>58155ca</code>.
</em></p>
