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
<hr/>
<p><em>
Generated with <code>gen-crd-api-reference-docs</code>
on git commit <code>2e776f6</code>.
</em></p>
