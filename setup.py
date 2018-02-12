from setuptools import setup
from pip.req import parse_requirements
from pip.download import PipSession


setup(
	name='ion',
	version='0.1',
	packages=['ion'],
	py_modules=['__main__'],
    install_requires=[str(ir.req) for ir in parse_requirements('requirements.txt', session=PipSession())],
	entry_points='''
	[console_scripts]
	ion=ion.__main__:main
	'''
)
