## Copyright 2017 Mhd Sulhan <ms@kilabit.info>. All rights reserved.
## Use of this source code is governed by a BSD-style
## license that can be found in the LICENSE file.

import logging

def getChanges(request, options=None):
    author = "buildbot"
    category = "mattermost"
    comments = "Build from mattermost"
    files = []
    properties = {}
    revlink = ""

    def firstOrNothing(value):
        """
        Small helper function to return the first value (if value is a list)
        or return the whole thing otherwise
        """
        if (isinstance(value, type([]))):
            return value[0]
        else:
            return value

    def do_build(cmd_args):
        repo = cmd_args[0]

        if len(cmd_args) >= 2:
            branch = cmd_args[1]
        else:
            branch = "master"

        ## which revision?
        if len(cmd_args) >= 3:
            revision = cmd_args[2]
        else:
            revision = ''

        project = "mattermost_build_%s_%s" % (repo, branch)

        ## which repository?
        repository = opt_repo_base_url + repo

        chdict = dict(
                author=author
            ,   branch=branch
            ,   category=category
            ,   comments=comments
            ,   files=files
            ,   project=project
            ,   properties=properties
            ,   repository=repository
            ,   revision=revision
            ,   revlink=revlink
        )

        return ([chdict], None)

    if options is None:
        raise ValueError('Missing options!')

    ## Get option values
    opt_token = options.get("token", None)
    if opt_token is None:
        raise ValueError("Missing option 'token'")

    opt_repo_base_url = options.get("repo_base_url", None)
    if opt_repo_base_url is None:
        raise ValueError("Missing option 'repo_base_url'")

    ## Parsing args
    args = request.args

    ##
    ## Check token payload.
    ##
    arg_token = firstOrNothing(args.get('token'))

    if opt_token != arg_token:
        raise ValueError('Invalid token!')

    ##
    ## Only allow specific channel to run the command.
    ##
    opt_channel = options.get('channel', None)
    arg_channel = firstOrNothing(args.get('channel_name'))

    if opt_channel != arg_channel:
        raise ValueError('Only specific channel can run the command %s %s' % (opt_channel, arg_channel))

    ## parsing command parameters
    command = firstOrNothing(args.get('command'))

    text = firstOrNothing(args.get('text'))
    cmd_args = text.split()

    if len(cmd_args) == 0:
        raise ValueError('Missing parameters...')

    if command != "/build":
        raise ValueError("Wrong command line!")

    return do_build(cmd_args)
