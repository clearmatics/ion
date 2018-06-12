# Copyright (c) 2018 Harry Roberts. All Rights Reserved.
# SPDX-License-Identifier: LGPL-3.0+

# Derived from https://gist.githubusercontent.com/HarryR/d2373f421c39353cd462/raw/c3f726455bbbc03fe791dda0dabbf4f73b5d2ec9/restclient.py
# Except the command-line interface has been removed (so it doesn't depend on Plugin and Host components)

__all__ = ('RestClient',)

try:
    from urllib.parse import quote_plus
except ImportError:
    from urllib import quote_plus

import requests

from .utils import require


class RestClient(object):
    __slots__ = ('_api', '_url', '_session')

    def __init__(self, url, api=None):
        require(url is not None, "Must provide REST API HTTP URL")
        self._url = url
        self._session = None
        if api is None:
            self._session = requests.Session()
        self._api = self if api is None else api

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

    def _request(self, method, kwargs):
        sess = self._api._session
        req = requests.Request(method=method,
                               url=self._url,
                               params=kwargs)
        resp = sess.send(req.prepare())
        resp.raise_for_status()
        return resp.json()

    def GET(self, **kwargs):
        return self._request('GET', kwargs)

    def POST(self, **kwargs):
        return self._request('POST', kwargs)

    def PUT(self, **kwargs):
        return self._request('PUT', kwargs)

    def DELETE(self, id=None, **kwargs):
        url = self._url
        if id is not None:
            url += "/" + str(id)
        sess = self._api._session
        req = requests.Request(method='DELETE',
                               url=url,
                               data=kwargs)
        resp = sess.send(req.prepare())
        resp.raise_for_status()
        return resp.json()
