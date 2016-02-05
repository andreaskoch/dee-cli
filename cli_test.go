package main

import (
	"github.com/andreaskoch/dee-ns"
	"net"
)

type testDNSEditorFactory struct {
	editor deens.DNSRecordEditor
	err    error
}

func (editorFactory testDNSEditorFactory) CreateDNSEditor() (deens.DNSRecordEditor, error) {
	return editorFactory.editor, editorFactory.err
}

type testDNSEditor struct {
	createSubdomainFunc func(domain, subDomainName string, timeToLive int, ip net.IP) error
	updateSubdomainFunc func(domain, subDomainName string, ip net.IP) error
	deleteSubdomainFunc func(domain, subDomainName string, recordType string) error
}

func (editor testDNSEditor) CreateSubdomain(domain, subDomainName string, timeToLive int, ip net.IP) error {
	return editor.createSubdomainFunc(domain, subDomainName, timeToLive, ip)
}

func (editor testDNSEditor) UpdateSubdomain(domain, subDomainName string, ip net.IP) error {
	return editor.updateSubdomainFunc(domain, subDomainName, ip)
}

func (editor testDNSEditor) DeleteSubdomain(domain, subDomainName string, recordType string) error {
	return editor.deleteSubdomainFunc(domain, subDomainName, recordType)
}

type testInfoProviderFactory struct {
	infoProvider deens.DNSInfoProvider
	err          error
}

func (factory testInfoProviderFactory) CreateInfoProvider() (deens.DNSInfoProvider, error) {
	return factory.infoProvider, factory.err
}
