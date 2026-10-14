# SAF-T (PT) Frequently Asked Questions

This document contains frequently asked questions about SAF-T (PT) as published by the Portuguese Tax Authority (Autoridade Tributária e Aduaneira).

## What is the SAF-T (PT) file?

SAF-T (PT) (Standard Audit File for Tax Purposes – Portuguese version) is a standardized file (in XML format) designed to allow easy export of a predefined set of accounting and billing records in a readable and common format, regardless of the program used, without affecting the internal structure of the program's database or its functionality.

The adoption of this model provides companies with a tool to meet the information requirements of inspection services, simplifying procedures and promoting the use of new technologies.

The file generation by information systems always covers a specific annual tax period, total or partial, from the beginning of that period until its end or until the generation date if earlier. Users are not allowed to choose which documents, series, or types of documents to export to SAF-T (PT).

## What are the advantages of SAF-T (PT) and who is it for?

The SAF-T (PT) file is designed to facilitate the electronic collection of relevant tax data from taxpayers by tax inspectors/auditors, as support for their tax returns and/or for the analysis of accounting records or others with tax relevance, as well as for compliance with declarative obligations that require it.

## Is there a tool to validate the consistency of the data exported in the SAF-T (PT) file?

SAF-T (PT) is a way to represent the data in the data repository. There is no official validation tool, so data consistency must be ensured by the application that created and exports it. Data export cannot alter or exclude existing records in the database – see paragraph a) of no. 1 of Annex I to Ordinance No. 302/2016, of December 2.

## Can I acquire an application that reads my application's data repository to generate SAF-T (PT) or integrate its data into another one for that purpose?

No. Accounting and/or billing applications or equivalent must themselves export the records from the databases they produce.

It is possible to integrate data from one application (integrated) into another (integrating) and export a SAF-T (PT) with data from both applications in the latter. However, the integrated application must always be able to export its own file - see paragraphs a) and f) of no. 1 of Annex I to Ordinance No. 302/2016, of December 2.

## In Annex I to Ordinance No. 302/2016, of December 2, there are "mandatory" fields and others without this mention. Does this mean that it is up to the user or software producer to define whether these "non-mandatory" fields are exported?

It is up to the software producer to generate the SAF-T (PT) file with the structure provided in Ordinance No. 302/2016, of December 2, for all mandatory fields, with the information existing in the application, and also for all other fields without this mention, belonging to the same structure, for which there is information in the data repository for their completion.

The user should be prevented from having any option about whether a particular operation or record should be included or excluded in the SAF-T (PT) file.

Not all fields are considered mandatory in the validation scheme, as there are some alternative choices and others whose generation may not make sense in a certain context. For example, if no discount has been given, it will not be relevant to require the completion of the "Line discount amount" (_SettlementAmount_) field.

Similarly, non-mandatory fields should not be generated when there is no information in the data repository for their completion. This procedure reduces the time spent generating the file as well as its size. Additionally, in SAF-T (PT), there cannot be empty fields under penalty of violation of the validation scheme [xsd file associated with SAF-T (PT) production].

## Is it possible to export SAF-T (PT) for periods shorter than the fiscal year?

SAF-T (PT), as a general rule, covers a complete fiscal year. Regarding the generation of SAF-T (PT) with accounting records, it must mandatorily include the complete fiscal year in a single file.

Exceptionally, and only for technical reasons associated with the size of commercial document tables, it is possible to generate the billing SAF-T (PT) file for complete monthly periods.

## Is it possible to change accounting software at a time not coinciding with the start of the fiscal year and start using the new program by performing only a balance migration?

The accounting SAF-T (PT) must be generated in a single file. The new accounting application must ensure the generation of SAF-T (PT) with the records made in the previous application, even if reported to the new references and nomenclatures, plus the records after its entry into operation.

The transition of records to a new accounting program at a time not coinciding with the start of the fiscal year cannot be done only with balance migration.

The replaced accounting application must also ensure the generation of accounting SAF-T (PT) until the moment of its discontinuation.

## Where can I find example files or validators to assess the SAF-T (PT) file structure?

There is an example SAF-T (PT) file in XML format available at [Portuguese Tax Authority's Software Certification page](https://info.portaldasfinancas.gov.pt/pt/apoio_contribuinte/Faturacao/Paginas/certificacao-de-software.aspx).

The structural validation tool for the file can be found at [SAF-T (PT) News page](https://info.portaldasfinancas.gov.pt/pt/apoio_contribuinte/SAFT_PT/Paginas/news-saf-t-pt.aspx) where you can also consult the data structure.
