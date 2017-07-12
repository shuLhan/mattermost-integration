## Copyright 2017 Mhd Sulhan <ms@kilabit.info>. All rights reserved.
## Use of this source code is governed by a BSD-style
## license that can be found in the LICENSE file.

## -*- python -*-
## vi: set tabstop=4 softtabstop=4 noexpandtab shiftwidth=4:
## vi: set comments=b\:##,fb\:- foldmarker=##[[,##]]:-

from twisted.python import log
from buildbot.process.results import SUCCESS
import json

"""
build parameter with properties:
{
	"buildrequestid": 11,
,	"complete_at": datetime.datetime(2017, 6, 1, 15, 31, 34, 851777, tzinfo=tzutc())
,	"complete": True
,	"builderid": 2
,	"buildid": 11
,	"workerid": 3
,	"builder": {
		"masterids": [2]
	,	"tags": []
	,	"description": u"Project description"
	,	"name": u"project-example"
	,	"builderid": 2
	}
,	"results": 2
,	"number": 9
,	"masterid": 2
,	"url": "http://172.32.1.xxx:8010/#builders/2/builds/9"
,	"buildrequest": {
		"buildrequestid": 11
	,	"complete": False
	,	"waited_for": False
	,	"claimed_at": datetime.datetime(2017, 6, 1, 15, 31, 34, tzinfo=tzutc())
	,	"results": -1
	,	"claimed": True
	,	"buildsetid": 11
	,	"complete_at": None
	,	"submitted_at": datetime.datetime(2017, 6, 1, 15, 31, 34, tzinfo=tzutc())
	,	"builderid": 2
	,	"claimed_by_masterid": 2
	,	"priority": 0
	}
,	"buildset": {
		"bsid": 11
	,	"complete_at": None
	,	"complete": False
	,	"sourcestamps": [{
			"project": u"mattermost_build_project-example_master"
		,	"codebase": u""
		,	"ssid": 4
		,	"branch": u"master"
		,	"repository": u"https://github.com/xxx/project-example"
		,	"patch": None
		,	"created_at": datetime.datetime(2017, 6, 1, 12, 54, 23, 75534, tzinfo=tzutc())
		,	"revision": u""
		}]
	,	"parent_buildid": None
	,	"results": -1
	,	"parent_relationship": None
	,	"reason": u"The SingleBranchScheduler scheduler named "mattermost_project-example" triggered this build"
	,	"external_idstring": None
	,	"submitted_at": 1496331094
	}
,	"state_string": u"failed "[ -d .git ] || git clone git@github.com/xxx/project-example .  git reset --hard HEAD ..." (failure) "cd doc ..." (failure)"
,	"started_at": datetime.datetime(2017, 6, 1, 15, 31, 34, 377513, tzinfo=tzutc())
,	"properties": {
		u"project": (u"mattermost_build_project-example_master", u"Build")
	,	u"branch": (u"master", u"Build")
	,	u"repository": (u"https://github.com/xxx/project-example", u"Build")
	,	u"buildername": (u"project-example", u"Builder")
	,	u"codebase": (u"", u"Build")
	,	u"slavename": (u"buildworker", u"Worker (deprecated)")
	,	u"workername": (u"buildworker", u"Worker")
	,	u"scheduler": (u"mattermost_project-example", u"Scheduler")
	,	u"builddir": (u"/data/buildworker/project-example", u"worker")
	,	u"buildnumber": (9, u"Build")
	,	u"revision": (u"", u"Build")
	}
}
"""

def formatter(build):
	mm_payload = dict(
		channel="ci"
	,   username="buildbot"
	,   text=""
	)

	repo = build["builder"]["name"]
	branch = build["properties"]["branch"][0]
	status = ":rocket:"
	build_url = build["url"]

	if build["complete"]:
		if build["results"] == SUCCESS:
			status = ":white_check_mark:"
		else:
			status = ":x:"

		mm_payload["text"] = "%s Finished %s at %s : %s\n" % (
			status, repo, branch, build_url)

	else:
		mm_payload["text"] = "%s Building %s at %s : %s\n" % (
			status, repo, branch, build_url)

	return mm_payload
