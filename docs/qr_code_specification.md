# QR Code Technical Specifications for Portuguese Invoicing

*Portaria n.º 195/2020*
Version 1.1

## Change History

| Date | Changes |
|------|----------|
| 2020-08-13 | Initial document creation |
| 2020-10-30 | - Updated QR message composition rules and field instructions<br>- Modified scope of code "L"<br>- Added new examples |

## Table of Contents

1. [Introduction](#1-introduction)
2. [Two-dimensional Barcode (QR Code)](#2-two-dimensional-barcode-qr-code)
3. [QR Code Message Composition](#3-qr-code-message-composition)
4. [QR Code Structure and Contents](#4-qr-code-structure-and-contents)
5. [QR Code Construction Examples](#5-qr-code-construction-examples)
   5.1. [Example 1 - Invoice](#51-example-1---invoice)
   5.2. [Example 2 - Simplified Invoice](#52-example-2---simplified-invoice)
   5.3. [Example 3 - Pro Forma Invoice](#53-example-3---pro-forma-invoice)
   5.4. [Example 4 - Transport Document](#54-example-4---transport-document)
   5.5. [Example 5 - Invoice with Foreign Tax Rate](#55-example-5---invoice-with-foreign-tax-rate)
   5.6. [Example 6 - Tax Correction Debit Note](#56-example-6---tax-correction-debit-note)
   5.7. [Example 7 - Invoice with VAT - Margin Scheme](#57-example-7---invoice-with-vat---margin-scheme)

## 1. Introduction

Decree-Law No. 28/2019 of February 15 consolidated and modernized invoicing regulations, introducing innovative aspects such as the unique document code and two-dimensional barcode (QR code). These changes aim to simplify invoice communication by individuals for determining IRS-deductible expenses while simultaneously enhancing control over taxpayer operations to combat informal economy, fraud, and tax evasion.

The two-dimensional barcode (QR code) must be included on invoices and other fiscally relevant documents issued by programs certified by the Portuguese Tax and Customs Authority (AT), as per Article 4 of Decree-Law No. 28/2019 of February 15.

This document defines the technical specifications for generating the two-dimensional barcode (QR code), as established in Ordinance No. 195/2020 of August 13.

Producers of AT-certified invoicing programs must ensure correct QR code generation according to these technical specifications and guarantee perfect readability within the document body, regardless of the medium presented to the customer.

For documents with multiple pages, the two-dimensional barcode (QR code) may appear on either the first or last page.

To support invoicing program producers in implementing QR codes in their programs, this document identifies the fields that must be included in the message, with their respective filling instructions and examples for each field.

QR code construction examples are also provided, covering various document types and different messages.

## 2. Two-dimensional Barcode (QR Code)

The QR code must comply with the following specifications:

a) Error Correction Level (ECC): "M"
b) Type: Byte
c) Points per module (Size): 2
d) Version: v=9 (minimum value)
e) Image dimensions: minimum 30×30 millimeters
f) Safety margin (Margin): 0.25 cm

## 3. QR Code Message Composition

When creating the message to be incorporated in the QR code, observe these rules:

a) Each field is formed by concatenating the "Code" column value from the table in section 4, ":" (colon), and the corresponding "Description" column value, without spaces.

b) The fields created this way must be concatenated in the order shown in the table in section 4, using "*" (asterisk) as separator.

c) Monetary fields must be represented in euros, using "." (period) as the decimal separator and always with two (2) decimal places, regardless of the number of decimal places (fewer or more) shown in the database and SAF-T (PT) file export. For documents issued in currency other than Euro, amounts must be converted to euros before presentation.

d) Fields marked with "+" are mandatory.

e) Fields marked with "++" are optional but should be created whenever information exists for them.

f) For optional fields, when information is absent, the field should not be created.

g) When creating fields with codes I1 to I8, J1 to J8, and K1 to K8, representing fiscal spaces for VAT purposes (e.g., PT, PTAC, and PT-MA), there must always be at least one fiscal space, up to a maximum of three simultaneous fiscal spaces, national or foreign.

h) When composing field with code I1, for documents issued without VAT rate indication that should appear in SAF-T (PT) tables 4.2, 4.3, or 4.4, it should be filled with "I1":"0".

i) No value in the "Description" column may exceed the maximum size defined in the table in section 4.

j) When composing field with code "S", when necessary, its elements should be concatenated with ";" (semicolon) without spaces.

k) For different VAT rates in the document (exempt, reduced, intermediate, and normal), the respective fields must show accumulated totals of tax base and VAT.

## 4. QR Code Structure and Contents

The QR code generation must follow the structure defined in the following table:

| Code | Description | Max Size | Instructions | Example | Required |
|------|-------------|----------|--------------|---------|----------|
| A | Issuer Tax ID | 9 | Fill with issuer's Tax ID without spaces and country prefix, according to SAF-T (PT) TaxRegistrationNumber field | A:123456789 | + |
| B | Customer Tax ID | 30 | Fill with customer's Tax ID without country prefix, according to SAF-T (PT) CustomerTaxID field. For "Final Consumer" documents, fill with 999999990 | B:999999990 | + |
| C | Customer Country | 12 | Fill according to the Country field in SAF-T (PT) customer table | C:PT | + |
| D | Document Type | 2 | Fill according to SAF-T (PT) typology - InvoiceType, MovementType, WorkType, or PaymentType fields | D:FT | + |
| E | Document Status | 1 | Fill according to SAF-T (PT) typology - InvoiceStatus, MovementStatus, WorkStatus, or PaymentStatus fields | E:N | + |
| F | Document Date | 8 | Use YYYYMMDD format. Corresponds to SAF-T (PT) InvoiceDate, MovementDate, WorkDate, or TransactionDate fields without hyphens | F:20191231 | + |
| G | Unique Document ID | 60 | Fill according to SAF-T (PT) typology - InvoiceNo, DocumentNumber, or PaymentRefNo fields | G:FT AB2019/0035 | + |
| H | ATCUD | 70 | Fill with unique document code according to SAF-T (PT) ATCUD fields | H:CSDF7T5H0035 | + |

[Table continues with all fields through S...]

## 5. QR Code Construction Examples

### 5.1. Example 1 - Invoice

Message with elements from the above table, including IBAN and VAT values for fiscal spaces PT, PT-AC, and PT-MA:

```
A:123456789*B:999999990*C:PT*D:FT*E:N*F:20191231*G:FT AB2019/0035*H:CSDF7T5H-0035*I1:PT*I2:12000.00*I3:15000.00*I4:900.00*I5:50000.00*I6:6500.00*I7:80000.00*I8:18400.00*J1:PT-AC*J2:10000.00*J3:25000.56*J4:1000.02*J5:75000.00*J6:6750.00*J7:100000.00*J8:18000.00*K1:PT-MA*K2:5000.00*K3:12500.00*K4:625.00*K5:25000.00*K6:3000.00*K7:40000.00*K8:8800.00*L:100.00*M:25.00*N:64000.02*O:513600.58*P:100.00*Q:kLp0*R:9999*S:TB;PT000000000000000000000000;513500.58
```

### 5.2. Example 2 - Simplified Invoice

Fields used in a simplified invoice:

| Code | Description | Example | Required |
|------|-------------|---------|----------|
| A | Issuer Tax ID | A:123456789 | + |
| B | Customer Tax ID | B:999999990 | + |
| C | Customer Country | C:PT | + |
| D | Document Type | D:FS | + |
| E | Document Status | E:N | + |
| F | Document Date | F:20190812 | + |
| G | Unique Document ID | G:FS CDVF/12345 | + |
| H | ATCUD | H:CDF7T5HD-12345 | + |
| I1 | Fiscal Space | I1:PT | + |
| I7 | Normal VAT Rate Tax Base | I7:0.65 | ++ |
| I8 | Normal VAT Rate Tax Amount | I8:0.15 | ++ |
| N | Total Tax Amount | N:0.15 | + |
| O | Document Total with Tax | O:0.80 | + |
| Q | Hash (4 characters) | Q:YhGV | + |
| R | Certificate Number | R:9999 | + |
| S | Other Information | S:NU;0.80 | ++ |

Generated message:
```
A:123456789*B:999999990*C:PT*D:FS*E:N*F:20190812*G:FS CDVF/12345*H:CDF7T5HD-12345*I1:PT*I7:0.65*I8:0.15*N:0.15*O:0.80*Q:YhGV*R:9999*S:NU;0.80
```

### 5.3. Example 3 - Pro Forma Invoice

Fields used in a pro forma invoice:

| Code | Description | Example | Required |
|------|-------------|---------|----------|
| A | Issuer Tax ID | A:500000000 | + |
| B | Customer Tax ID | B:123456789 | + |
| C | Customer Country | C:PT | + |
| D | Document Type | D:PF | + |
| E | Document Status | E:N | + |
| F | Document Date | F:20190123 | + |
| G | Unique Document ID | G:PF G2019CB/145789 | + |
| H | ATCUD | H:HB6FT7RV-145789 | + |
| I1 | Fiscal Space | I1:PT | + |
| I2 | Exempt Tax Base | I2:12345.34 | ++ |
| I3 | Reduced VAT Rate Tax Base | I3:12532.65 | ++ |
| I4 | Reduced VAT Rate Tax Amount | I4:751.96 | ++ |
| I5 | Intermediate VAT Rate Tax Base | I5:52789.00 | ++ |
| I6 | Intermediate VAT Rate Tax Amount | I6:6862.57 | ++ |
| I7 | Normal VAT Rate Tax Base | I7:32425.69 | ++ |
| I8 | Normal VAT Rate Tax Amount | I8:7457.91 | ++ |
| N | Total Tax Amount | N:15072.44 | + |
| O | Document Total with Tax | O:125165.12 | + |
| Q | Hash (4 characters) | Q:r/fY | + |
| R | Certificate Number | R:9999 | + |

Generated message:
```
A:500000000*B:123456789*C:PT*D:PF*E:N*F:20190123*G:PF G2019CB/145789*H:HB6FT7RV-145789*I1:PT*I2:12345.34*I3:12532.65*I4:751.96*I5:52789.00*I6:6862.57*I7:32425.69*I8:7457.91*N:15072.44*O:125165.12*Q:r/fY*R:9999
```

### 5.4. Example 4 - Transport Document

Fields used in a non-valued transport guide:

| Code | Description | Example | Required |
|------|-------------|---------|----------|
| A | Issuer Tax ID | A:500000000 | + |
| B | Customer Tax ID | B:123456789 | + |
| C | Customer Country | C:PT | + |
| D | Document Type | D:GT | + |
| E | Document Status | E:N | + |
| F | Document Date | F:20190720 | + |
| G | Unique Document ID | G:GT G234CB/50987 | + |
| H | ATCUD | H:GTVX4Y8B-50987 | + |
| I1 | Fiscal Space | I1:0 | + |
| N | Total Tax Amount (*) | N:0.00 | + |
| O | Document Total with Tax (*) | O:0.00 | + |
| Q | Hash (4 characters) | Q:5ulg | + |
| R | Certificate Number | R:9999 | + |

(*) Even if these values exist in the program's database, if the document doesn't display the values, they should not appear in the QR code. The same procedure should be followed for documents that may be presented to the customer and allow verification of goods or services that are valued in the database when the values do not appear on the documents.

Generated message:
```
A:500000000*B:123456789*C:PT*D:GT*E:N*F:20190720*G:GT G234CB/50987*H:GTVX4Y8B-50987*I1:0*N:0.00*O:0.00*Q:5ulg*R:9999
```

### 5.5. Example 5 - Invoice with Foreign Tax Rate

Fields used in an invoice issued with French VAT rate, with a value of 100 euros subject to standard rate VAT (20%):

| Code | Description | Example | Required |
|------|-------------|---------|----------|
| A | Issuer Tax ID | A:123456789 | + |
| B | Customer Tax ID | B:4443332215 | + |
| C | Customer Country | C:FR | + |
| D | Document Type | D:FT | + |
| E | Document Status | E:N | + |
| F | Document Date | F:20190526 | + |
| G | Unique Document ID | G:ABC BNH/4561 | + |
| H | ATCUD | H:DK5ZJ2HN-4561 | + |
| I1 | Fiscal Space | I1:FR | + |
| I7 | Normal VAT Rate Tax Base | I7:100.00 | ++ |
| I8 | Normal VAT Rate Tax Amount | I8:20.00 | ++ |
| N | Total Tax Amount | N:20.00 | + |
| O | Document Total with Tax | O:120.00 | + |
| Q | Hash (4 characters) | Q:YJRE | + |
| R | Certificate Number | R:9999 | + |

Generated message:
```
A:123456789*B:4443332215*C:FR*D:FT*E:N*F:20190526*G:ABC BNH/4561*H:DK5ZJ2HN-4561*I1:FR*I7:100.00*I8:20.00*N:20.00*O:120.00*Q:YJRE*R:9999
```

### 5.6. Example 6 - Tax Correction Debit Note

Fields used in a debit note issued only to correct insufficiently charged tax in an invoice, amounting to €26.50:

| Code | Description | Example | Required |
|------|-------------|---------|----------|
| A | Issuer Tax ID | A:123456789 | + |
| B | Customer Tax ID | B:500000000 | + |
| C | Customer Country | C:PT | + |
| D | Document Type | D:ND | + |
| E | Document Status | E:N | + |
| F | Document Date | F:20190216 | + |
| G | Unique Document ID | G:M1F KLG/6145 | + |
| H | ATCUD | H:RQD8L6DG-6145 | + |
| I1 | Fiscal Space | I1:PT-MA | + |
| I6 | Intermediate VAT Rate Tax Amount | I6:26.50 | ++ |
| N | Total Tax Amount | N:26.50 | + |
| O | Document Total with Tax | O:26.50 | + |
| Q | Hash (4 characters) | Q:h1rB | + |
| R | Certificate Number | R:9999 | + |

Generated message:
```
A:123456789*B:500000000*C:PT*D:ND*E:N*F:20190216*G:M1F KLG/6145*H:RQD8L6DG-6145*I1:PT-MA*I6:26.50*N:26.50*O:26.50*Q:h1rB*R:9999
```

### 5.7. Example 7 - Invoice with VAT - Margin Scheme

Fields used in an invoice issued under a margin scheme VAT regime (€1,000.00 including tax), also containing a transaction subject to the general VAT regime (€59.00 including tax):

| Code | Description | Example | Required |
|------|-------------|---------|----------|
| A | Issuer Tax ID | A:500000000 | + |
| B | Customer Tax ID | B:123456789 | + |
| C | Customer Country | C:PT | + |
| D | Document Type | D:FT | + |
| E | Document Status | E:N | + |
| F | Document Date | F:20191124 | + |
| G | Unique Document ID | G:NF 19A/789145 | + |
| H | ATCUD | H:JL9DS4TT-789145 | + |
| I1 | Fiscal Space | I1:PT-AC | + |
| I7 | Normal VAT Rate Tax Base | I7:50.00 | ++ |
| I8 | Normal VAT Rate Tax Amount | I8:9.00 | ++ |
| L | Non-subject/non-taxable VAT/other situations | L:1000.00 | ++ |
| N | Total Tax Amount (*) | N:9.00 | + |
| O | Document Total with Tax | O:1059.00 | + |
| Q | Hash (4 characters) | Q:d8/K | + |
| R | Certificate Number | R:9999 | + |

(*) If the program can determine the VAT amount underlying the margin scheme transaction and it exists in its data repository, it must be exported to the SAF-T (PT) file. However, since legally issued documents cannot itemize VAT, it also cannot appear in the QR code.

Generated message:
```
A:500000000*B:123456789*C:PT*D:FT*E:N*F:20191124*G:NF 19A/789145*H:JL9DS4TT-789145*I1:PT-AC*I7:50.00*I8:9.00*L:1000.00*N:9.00*O:1059.00*Q:d8/K*R:9999
```
