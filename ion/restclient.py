# Copyright (c) 2018 Harry Roberts. All Rights Reserved.
# SPDX-License-Identifier: LGPL-3.0+

"""
Derived from https://gist.githubusercontent.com/HarryR/d2373f421c39353cd462/raw/c3f726455bbbc03fe791dda0dabbf4f73b5d2ec9/restclient.py

Except the command-line interface has been removed (so it doesn't depend on
Plugin and Host components). I didn't realise how broken the original Gist was,
this is less broken and much cleaner....

This class provides a friendly Pythonic interface to a REST-ish HTTP JSON API.
Form parameters are posted as normal form encoded data, Response is returned as
JSON.

e.g.

    x = RestClient('http://example.com/')
    x('test').POST(abc=123) # POST /test abc=123
    x.test.derp.GET()       # GET /test/derp
    x.test.GET()()          # GET /test

see... it's nice, and predictable, and Pythonic, and flexible, etc...
"""

__all__ = ('RestClient',)

try:
    from urllib.parse import quote_plus
except ImportError:
    from urllib import quote_plus

import requests


class RestClient(object):
    __slots__ = ('_api', '_url', '_session')

    def __init__(self, url, api=None):
        assert url is not None
        self._url = url
        self._session = None
        if api is None:
            self._session = requests.Session()
        self._api = self if api is None else api
        assert isinstance(self._api, RestClient)

    def __getattr__(self, name):
        if name[0] == '_':
            raise AttributeError
        return RestClient(self._url + '/' + name, self._api)

    def __call__(self, name=None):
        if name is None:
            return self
        if name[0] == '_':
            raise AttributeError
        return RestClient(self._url + '/' + quote_plus(name), self._api)

    def _request(self, method, params=None, data=None):
        sess = self._api._session
        req = requests.Request(method=method,
                               url=self._url,
                               params=params,
                               data=data)
        resp = sess.send(req.prepare())
        try:
            data = resp.json()
            error = data.get('_error')
            if error:
                raise RuntimeError(error)
            return data
        except ValueError:
            resp.raise_for_status()

    def GET(self, **kwargs):
        return self._request('GET', params=kwargs)

    def POST(self, **kwargs):
        return self._request('POST', data=kwargs)

    def PUT(self, **kwargs):
        return self._request('PUT', data=kwargs)

    def DELETE(self, id=None, **kwargs):
        return self._request('DELETE', data=kwargs)
