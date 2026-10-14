# Legislation

## Decree - Ordinance No. 302/2016, of December 2

**Status:** in force

**Summary:** Amends the data structure of the file referred to in Ordinance No. 321-A/2007, of March 26, and creates the taxonomies to be used in filling out the SAF-T (PT) file.

**Publication:** Official Gazette No. 231/2016, Series I of 2016-12-02, pages 4273 - 4379

**Associated legislation:** Ordinance No. 321-A/2007 - 26/03

**History of amendments:** Declaration of Rectification No. 2-A/2017 - 02/02

[See - original DR](#)

## FINANCE

### Ordinance No. 302/2016, of December 2

Ordinance No. 321-A/2007, of March 26, approved a standardized tax audit file format for data export, the so-called SAF-T (PT), which has proven to be an excellent tool for obtaining information by inspection services and whose data structure has been adapted according to accounting or tax changes.

The evolution in the SAF-T (PT) file data structure has focused mainly on improving the quality of billing information. The experience of using SAF-T (PT) has shown that the current structure is insufficient for a complete understanding and control of accounting information, due to the flexibility in the use of accounts by different entities.

In this perspective, it is important to adjust the structure of the SAF-T (PT) file with the creation of taxonomies, that is, correspondence tables that allow the characterization of accounts according to the accounting standards used by different taxpayers, allowing to simplify the filling of Annexes A and I of the IES.

Therefore:
The Government orders, through the Secretary of State for Tax Affairs, under the provisions of paragraph 8 of article 123 of the IRC Code, approved by Decree-Law No. 442-B/88, of November 30, the following:

### Article 1

**Object**
This ordinance:
a) Amends Ordinance No. 321-A/2007, of March 26;
b) Amends the data structure of the file referred to in paragraph 1 of Ordinance No. 321-A/2007, of March 26;
c) Creates the taxonomies to be used in filling out the fields duly marked in the SAF-T (PT) file data structure.

### Article 2

**Amendment to Ordinance No. 321-A/2007, of March 26**

Paragraph 1 of Ordinance No. 321-A/2007, of March 26, shall read as follows:
"1. The taxpayers referred to in paragraph 1 of article 123 of the Corporate Income Tax Code are required to produce a file, according to the data structure in the annex, whenever requested by the Tax and Customs Inspection services, within the scope of their competencies, or for compliance with declarative obligations that require it."

### Article 3

**Data Structure**

The file referred to in paragraph 1 of Ordinance No. 321-A/2007, of March 26, with the last amendment introduced by Ordinance No. 274/2013, of August 21, shall have the data structure contained in Annex I to this ordinance, which forms an integral part thereof.

### Article 4

**Taxonomies**
The filling out of fields relating to the account classification code in the data structure of the file referred to in paragraph 1 of Ordinance No. 321-A/2007, of March 26, must be done with reference to the taxonomies contained in Annexes II and III to this ordinance, which form an integral part thereof, depending on whether the account classification reference used by each entity corresponds to SNC Base or International Accounting Standards (Annex II) or SNC Microentities (Annex III), respectively.

### Article 5

**Entry into force**

This ordinance enters into force on January 1, 2017, for the fiscal years 2017 and following, with the exception of the data structure referred to in Article 3, which enters into force on July 1, 2017.

The Secretary of State for Tax Affairs, Fernando António Portela Rocha de Andrade, on November 29, 2016.

## Annex I

### Data Structure

(referred to in Article 3 of this ordinance)

#### 1 - GENERIC ASPECTS

a) Accounting and/or billing applications (including those that issue transport documents provided for in Decree-Law No. 147/2003, of July 11, and other documents that may be presented to the customer for checking goods or services) and those that issue receipts must, by themselves, export the predefined set of database records they produce, in a readable and common format, with the data structure and respective restrictions provided in the validation scheme, without affecting the internal structure of the program's database or its functionality.

b) The SAF-T (PT) file must be generated in a standardized format, in XML language, respecting not only the validation scheme in the xsd format file that is available at [http://www.portaldasfinancas.gov.pt](http://www.portaldasfinancas.gov.pt), but also the content specified in this ordinance.

c) The generation of the SAF-T (PT) file by information systems must always be carried out for a specific annual taxation period, total or partial, from the beginning of that period until its end or until the generation date if earlier.

d) In the "Mandatory" column, the symbol "*" corresponds to a mandatory field and the symbol "**" corresponds to a field of alternative choice or dependent on conditions for its mandatory nature, and the remaining fields must also be filled in as long as the information exists in the application's data repository, which can be verified in the respective documentation.

e) The SAF-T (PT) file relating to accounting must be unique for the period to which it refers. In the case of SAF-T (PT) files relating to billing, it is possible to generate one for each establishment, if independent of the billing system adopted at the central level. If the billing system in the establishments is centralized, a single file must be provided.

f) Billing applications, even if used by third parties for issuing documents in the name and on behalf of taxpayers, and integrated accounting and billing applications, including those that issue transport documents provided for in Decree-Law No. 147/2003, of July 11, and other documents that may be presented to the customer for checking goods or services, must generate a file with all mandatory tables.

In the case of accounting and billing applications generating independent files, these must contain common information and specific information.

**Common information tables for both files:**

1. Header;
2. Customer table;
3. Tax table; and
4. Payment receipts issued (Payments), when it must exist.

**Specific information tables:**

For accounting applications:

1. Chart of accounts table (GeneralLedgerAccounts);
2. Supplier table; and
3. Accounting movements (GeneralLedgerEntries).

For billing applications, even if used by third parties for issuing documents in the name and on behalf of taxpayers, including those that issue transport documents and other documents that may be presented to the customer for checking goods or services, when they must exist:

1. Supplier table;
2. Products/services table (Product);
3. Commercial documents to customers (SalesInvoices);
4. Goods movement documents (MovementOfGoods); and
5. Goods or service provision checking documents (WorkingDocuments).

g) In the case of self-billing, the generation of the corresponding file is the responsibility of the actual issuer (the customer who self-bills), who must make it available whenever their supplier is required to provide the SAF-T (PT) file.

In this case, the issuer must provide a file with the data from the tables used, namely:

- Table 1. - Header with fields 1.2 to 1.8 with the supplier's data, where field 1.4 - Tax Accounting Basis (TaxAccountingBasis) must be filled in with code "S";
- Table 2.2. - Customers (Customer) with the issuer's data;
- Table 2.4. - Products/services (Product) with the records of products/services mentioned in the documents; and
- Table 4.1. - Commercial documents to customers (SalesInvoices) with documents relating to the supplier where field 4.1.4.3.1 - Current document status (InvoiceStatus) is filled in with code "S".

h) When billing is carried out by third parties, in the name and on behalf of the taxpayer, the file must contain the tables indicated in paragraph f); in Table 1 - Header, field 1.4 - Tax Accounting Basis (TaxAccountingBasis) must be filled in with "E", and in field 1.18 - Additional comments (HeaderComment) the actual issuer (TIN and name) and the location where the respective databases are must be identified.

i) When the requested files pose extraction problems, due to the size of the commercial document tables (SourceDocuments): Tables 4.1. to 4.4., which must exist, it is possible to subdivide this extraction by complete monthly periods, included in the request, complemented, if applicable, with the last period less than one month.

j) Dates are indicated in the date format: "YYYY-MM-DD" and time records in the date and time format: "YYYY-MM-DDThh:mm:ss", without including the time zone and milliseconds.

k) In the "Format" column, the dimension of text type fields is indicated in number of characters.

l) In the export of values, negative values are not accepted. When there are accounting balances or transactions in commercial documents with negative values, their representation must be adjusted to debit or credit and the remaining amount fields or values must be exported in absolute value.

#### 2 - DATA STRUCTURE

1. Header.

This table contains general information relating to the taxpayer to whom the SAF-T (PT) refers.

| Field Index | Mandatory | Field Name | Technical Notes | Format (validated in xsd) |
|:----------:|:---------:|:-----------|:----------------|:------------------------:|
| 1.1. | * | Audit file version (AuditFileVersion) | The XML schema version to be used will be the one available at [http://www.portaldasfinancas.gov.pt](http://www.portaldasfinancas.gov.pt) | Text 10 |
| 1.2. | * | Company commercial registration ID (CompanyID) | Obtained by concatenating the commercial registry office with the commercial registration number, separated by a space character. In cases where there is no commercial registration, the TIN must be indicated. | Text 50 |
| 1.3. | * | Company tax registration number (TaxRegistrationNumber) | Fill in with the Portuguese TIN without spaces and without any country prefix. | Integer 9 |
| 1.4. | * | Tax accounting basis (TaxAccountingBasis) | Must be filled in with the type of program providing the applicable data (including transport documents, checking documents and receipts issued, when they exist): "C" - Accounting; "E" - Billing issued by third parties; "F" - Billing; "I" - Accounting integrated with billing; "P" - Partial billing; "R" - Receipts (a); "S" - Self-billing; "T" - Transport documents (a). (a) This type of program should be indicated if it issues only this type of document. Otherwise, type "C", "F" or "I" should be used. | Text 1 |
| 1.5. | * | Company name (CompanyName) | Company's corporate name or taxpayer's name. | Text 100 |
| 1.6. | | Business name (BusinessName) | Taxpayer's business name. | Text 60 |
| 1.7. | * | Company address (CompanyAddress) | | N/A |
| 1.7.1. | | Building number (BuildingNumber) | | Text 10 |
| 1.7.2. | | Street name (StreetName) | | Text 200 |
| 1.7.3. | * | Address detail (AddressDetail) | Must include street name, building number and floor, if applicable. | Text 210 |
| 1.7.4. | * | City (City) | | Text 50 |
| 1.7.5. | * | Postal code (PostalCode) | | Text 8 |
| 1.7.6. | | Region (Region) | | Text 50 |
| 1.7.7. | * | Country (Country) | Fill in with "PT". | Text 2 |
| 1.8. | * | Fiscal year (FiscalYear) | Use IRC Code rules in case of accounting periods not coinciding with the calendar year. (Example: taxation period from 01-10-2012 to 30-09-2013 corresponds to FiscalYear = 2012). | Integer 4 |
| 1.9. | * | Start date of file period (StartDate) | | Date |
| 1.10. | * | End date of file period (EndDate) | | Date |
| 1.11. | * | Currency code (CurrencyCode) | Identifies the standard currency to be used in monetary type fields in the file. Fill in with "EUR". | Text 3 |
| 1.12. | * | Date created (DateCreated) | Creation date of the SAF-T (PT) XML file. | Date |
| 1.13. | * | Tax entity (TaxEntity) | In the case of billing file, it must be specified which establishment the produced file refers to, if applicable. Otherwise, it must be filled in with the specification "Global". In the case of accounting or integrated file, this field must be filled in with the specification "Sede". | Text 20 |
| 1.14. | * | Software producer tax ID (ProductCompanyTaxID) | Fill in with the TIN of the software producing entity. | Text 20 |
| 1.15. | * | Software certificate number (SoftwareCertificateNumber) | Certificate number assigned to the software producing entity, according to Ordinance No. 363/2010, of June 23. If not applicable, must be filled in with "0" (zero). | Integer |
| 1.16. | * | Product ID (ProductID) | Name of the application that generates the SAF-T (PT). The commercial name of the software and the company producing it must be indicated in the format "Application name/Software producer company name". | Text 255 |
| 1.17. | * | Product version (ProductVersion) | The version of the application producing the file must be indicated. | Text 30 |
| 1.18. | | Additional comments (HeaderComment) | | Text 255 |
| 1.19. | | Telephone (Telephone) | | Text 20 |
| 1.20. | | Fax (Fax) | | Text 20 |
| 1.21. | | Company email (Email) | | Text 254 |
| 1.22. | | Company website (Website) | | Text 60 |

2. Master Files (MasterFiles):

The master tables 2.1; 2.2; 2.3; 2.4 and 2.5 of this data structure are mandatory under the conditions referred to in paragraphs f), g), h) and i) of point 1 of this annex, when applicable.

#### 2.1. - Chart of accounts table (GeneralLedgerAccounts)

The chart of accounts table to be exported is the one provided by the accounting standardization system and other legal provisions in force for the respective sector of activity. Records of account classes should not be exported.

In the case of aggregating accounts that contain sub-accounts with debit balances and sub-accounts with credit balances, these debit and credit balances must be shown in the aggregating account.

| Field Index | Mandatory | Field Name | Technical Notes | Format (validated in xsd) |
|:----------:|:---------:|:-----------|:----------------|:------------------------:|
| 2.1.1. | * | Account classification reference (TaxonomyReference) | The account classification reference used by the entity (chart of accounts) must be indicated: "S" - SNC base - Taxonomy S; "N" - International Accounting Standards - Taxonomy S; "M" - SNC microentities - Taxonomy M; "O" - Other accounting references, whose taxonomy is not codified. | Text 1 |
| 2.1.2. | * | Account table record (Account) | | N/A |
| 2.1.2.1. | * | Account code (AccountID) | All accounts must be included in the file, including their respective integrating accounts, up to the General Ledger accounts. | Text 30 |
| 2.1.2.2. | * | Account description (AccountDescription) | | Text 100 |
| 2.1.2.3. | * | Opening debit balance (OpeningDebitBalance) | The opening debit balance will always be that of the beginning of the taxation period. | Monetary |
| 2.1.2.4. | * | Opening credit balance (OpeningCreditBalance) | The opening credit balance will always be that of the beginning of the taxation period. | Monetary |
| 2.1.2.5. | * | Closing debit balance (ClosingDebitBalance) | The closing debit balance will be that of the end of the taxation period or the generation date, if earlier. | Monetary |
| 2.1.2.6. | * | Closing credit balance (ClosingCreditBalance) | The closing credit balance will be that of the end of the taxation period or the generation date, if earlier. | Monetary |
| 2.1.2.7. | * | Account category and type (GroupingCategory) | The type and category of the account must be indicated: "GR" - First-degree account of general accounting; "GA" - Aggregating or integrating account of general accounting; "GM" - Movement account of general accounting; "AR" - First-degree account of analytical accounting; "AA" - Aggregating or integrating account of analytical accounting; and "AM" - Movement account of analytical accounting. | Text 2 |
| 2.1.2.8. | ** | Account hierarchy (GroupingCode) | Except for first-degree accounts, the respective aggregating account of the immediately higher degree must be indicated, using for this purpose the exact structure that appears in the corresponding field 2.1.2.1. - Account code (AccountID). | Text 30 |
| 2.1.2.9. | ** | Account classification code (TaxonomyCode) | Must indicate one of the Account classification codes, according to the indicated TaxonomyReference: - TaxonomyReference=S, TaxonomyCode of Taxonomy S; - TaxonomyReference=N, TaxonomyCode of Taxonomy S; - TaxonomyReference=M, TaxonomyCode of Taxonomy M; - TaxonomyReference=O, TaxonomyCode must always be equal to "1". This field is mandatory whenever GroupingCategory is equal to "GM". | Integer |

#### 2.2. - Customers table (Customer)

This table must contain all records moved in the taxation period in the respective customer file, as well as those that are implicit in the movements and do not appear in the respective file. If, for example, there is an invoice with only the tax identification number or customer name registered, which does not appear in the application's customer file, it must be exported as a customer in the SAF-T (PT).

| Field Index | Mandatory | Field Name | Technical Notes | Format (validated in xsd) |
|:----------:|:---------:|:-----------|:----------------|:------------------------:|
| 2.2.1. | * | Unique customer identifier (CustomerID) | In the customer list, there cannot be more than one record with the same CustomerID. For the case of final consumers, a generic customer must be created with the designation "Final consumer". | Text 30 |
| 2.2.2. | * | Account code (AccountID) | The respective customer current account in the chart of accounts must be indicated, if defined. Otherwise, it must be filled in with the designation "Unknown". | Text 30 |
| 2.2.3. | * | Customer tax identification number (CustomerTaxID) | Must be indicated without the country prefix. The generic customer, corresponding to the designated "Final consumer", must be identified with TIN "999999990". | Text 30 |
| 2.2.4. | * | Company name (CompanyName) | The generic customer must be identified with the designation "Final consumer". | Text 100 |
| 2.2.5. | | Company contact name (Contact) | | Text 50 |
| 2.2.6. | * | Billing address (BillingAddress) | Corresponds to the address of the headquarters or permanent establishment in national territory. | N/A |
| 2.2.6.1. | | Building number (BuildingNumber) | | Text 10 |
| 2.2.6.2. | | Street name (StreetName) | | Text 200 |
| 2.2.6.3. | * | Address detail (AddressDetail) | Must include street name, building number and floor, if applicable. Must be filled in with the designation "Unknown" in the following situations: - Non-integrated systems, if the information is not known; - Operations carried out with "Final consumer". | Text 210 |
| 2.2.6.4. | * | City (City) | Must be filled in with the designation "Unknown" in the following situations: - Non-integrated systems, if the information is not known; - Operations carried out with "Final consumer". | Text 50 |
| 2.2.6.5. | * | Postal code (PostalCode) | Must be filled in with the designation "Unknown" in the following situations: - Non-integrated systems, if the information is not known; - Operations carried out with "Final consumer". | Text 20 |
| 2.2.6.6. | | Region (Region) | | Text 50 |
| 2.2.6.7. | * | Country (Country) | If known, must be filled in according to ISO 3166-1-alpha-2 standard. Must be filled in with the designation "Unknown" in the following situations: - Non-integrated systems, if the information is not known; - Operations carried out with "Final consumer". | Text 12 |
| 2.2.7. | | Shipping address (ShipToAddress) | If there is a need to make more than one reference, this structure can be generated as many times as necessary. | N/A |
| 2.2.7.1. | | Building number (BuildingNumber) | | Text 10 |
| 2.2.7.2. | | Street name (StreetName) | | Text 200 |
| 2.2.7.3. | * | Address detail (AddressDetail) | Must include street name, building number and floor, if applicable. Must be filled in with the designation "Unknown" in the following situations: - Non-integrated systems, if the information is not known; - Operations carried out with "Final consumer". | Text 210 |
| 2.2.7.4. | * | City (City) | Must be filled in with the designation "Unknown" in the following situations: - Non-integrated systems, if the information is not known; - Operations carried out with "Final consumer". | Text 50 |
| 2.2.7.5. | * | Postal code (PostalCode) | Must be filled in with the designation "Unknown" in the following situations: - Non-integrated systems, if the information is not known; - Operations carried out with "Final consumer". | Text 20 |
| 2.2.7.6. | | Region (Region) | | Text 50 |
| 2.2.7.7. | * | Country (Country) | Must be filled in according to ISO 3166-1-alpha-2 standard. Must be filled in with the designation "Unknown" in the following situations: - Non-integrated systems, if the information is not known; - Operations carried out with "Final consumer". | Text 12 |
| 2.2.8. | | Telephone (Telephone) | | Text 20 |
| 2.2.9. | | Fax (Fax) | | Text 20 |
| 2.2.10. | | Company email (Email) | | Text 254 |
| 2.2.11. | | Company website (Website) | | Text 60 |
| 2.2.12. | * | Self-billing indicator (SelfBillingIndicator) | Indicator of the existence of a self-billing agreement between the customer and the supplier. Must be filled in with "1" if there is an agreement and with "0" (zero) otherwise. | Integer |

#### 2.3. - Suppliers table (Supplier)

This table must contain all records moved in the taxation period in the respective database.

| Field Index | Mandatory | Field Name | Technical Notes | Format (validated in xsd) |
|:----------:|:---------:|:-----------|:----------------|:------------------------:|
| 2.3.1. | * | Unique supplier identifier (SupplierID) | In the supplier list, there cannot be more than one record with the same SupplierID. | Text 30 |
| 2.3.2. | * | Account code (AccountID) | The respective supplier current account in the chart of accounts must be indicated, if defined. Otherwise, it must be filled in with the designation "Unknown". | Text 30 |
| 2.3.3. | * | Supplier tax identification number (SupplierTaxID) | Must be indicated without the country prefix. | Text 30 |
| 2.3.4. | | Company name (CompanyName) | | Text 100 |
| 2.3.5. | | Company contact name (Contact) | | Text 50 |
| 2.3.6. | | Billing address (BillingAddress) | Corresponds to the address of the headquarters or permanent establishment in national territory. | N/A |
| 2.3.6.1. | | Building number (BuildingNumber) | | Text 10 |
| 2.3.6.2. | | Street name (StreetName) | | Text 200 |
| 2.3.6.3. | | Address detail (AddressDetail) | Must include street name, building number and floor, if applicable. | Text 210 |
| 2.3.6.4. | | City (City) | | Text 50 |
| 2.3.6.5. | | Postal code (PostalCode) | | Text 20 |
| 2.3.6.6. | | Region (Region) | | Text 50 |
| 2.3.6.7. | | Country (Country) | Must be filled in according to ISO 3166-1-alpha-2 standard. | Text 2 |
| 2.3.7. | | Shipping from address (ShipFromAddress) | If there is a need to make more than one reference, this structure can be generated as many times as necessary. | N/A |
| 2.3.7.1. | | Building number (BuildingNumber) | | Text 10 |
| 2.3.7.2. | | Street name (StreetName) | | Text 200 |
| 2.3.7.3. | | Address detail (AddressDetail) | Must include street name, building number and floor, if applicable. | Text 210 |
| 2.3.7.4. | | City (City) | | Text 50 |
| 2.3.7.5. | | Postal code (PostalCode) | | Text 20 |
| 2.3.7.6. | | Region (Region) | | Text 50 |
| 2.3.7.7. | | Country (Country) | Must be filled in according to ISO 3166-1-alpha-2 standard. | Text 2 |
| 2.3.8. | | Telephone (Telephone) | | Text 20 |
| 2.3.9. | | Fax (Fax) | | Text 20 |
| 2.3.10. | | Company email (Email) | | Text 254 |
| 2.3.11. | | Company website (Website) | | Text 60 |
| 2.3.12. | * | Self-billing indicator (SelfBillingIndicator) | Indicator of the existence of a self-billing agreement between the customer and the supplier. Must be filled in with "1" if there is an agreement and with "0" (zero) otherwise. | Integer |

#### 2.4. Products/services table (Product)

This table must contain the catalog of products and types of services used in the billing system that were subject to movement and also the records that are implicit in the movements and that do not exist in the application's Products/Services table.

If, for example, there is an invoice with a shipping line that does not appear in the application's items file, this must be exported and represented as a product in the SAF-T (PT).

Taxes, fees, eco-taxes and parafiscal charges that appear mentioned in the invoice and that contribute or not to the VAT or Stamp Duty tax base must also be included - except VAT and Stamp Duty, which should be reflected in 2.5. - Tax table (TaxTable).

| Field Index | Mandatory | Field Name | Technical Notes | Format (validated in xsd) |
|:----------:|:---------:|:-----------|:----------------|:------------------------:|
| 2.4.1. | * | Product or service indicator (ProductType) | Must be filled in with: "P" - Products; "S" - Services; "O" - Others (Ex: charged shipping, received advances or asset disposal); "E" - Special Consumption Taxes - (ex.: IABA, ISP, IT); "I" - Taxes, fees and parafiscal charges - except VAT and Stamp Duty which should be reflected in table 2.5 - Tax table (TaxTable) and Special Consumption Taxes, which should be filled in with code "E". | Text 1 |
| 2.4.2. | * | Product or service identifier (ProductCode) | Unique code of the product in the product list. | Text 60 |
| 2.4.3. | | Product or service family (ProductGroup) | | Text 50 |
| 2.4.4. | * | Product or service description (ProductDescription) | Must correspond to the usual denomination of the goods or services provided, with specification of the elements necessary to determine the applicable tax rate. | Text 200 |
| 2.4.5. | * | Product code (ProductNumberCode) | The EAN code (barcode) of the product must be used. When this does not exist, fill in with the value of field 2.4.2. Product or service identifier (ProductCode). | Text 60 |
| 2.4.6. | | Customs details (CustomsDetails) | | N/A |
| 2.4.6.1. | | Combined Nomenclature Code (CNCode) | Fill in with the European Union combined nomenclature code. If there is a need to make more than one reference, this field can be generated as many times as necessary. | Text 8 |
| 2.4.6.2. | | UN Number (UNNumber) | Fill in with the UN number for dangerous products. If there is a need to make more than one reference, this field can be generated as many times as necessary. | Text 4 |

#### 2.5. Tax table (TaxTable)

In this table are recorded the VAT tax regimes, practiced in each fiscal space and the stamp duty items to be assessed, applicable in the document lines, recorded in table 4. - Commercial documents (SourceDocuments).

| Field Index | Mandatory | Field Name | Technical Notes | Format (validated in xsd) |
|:----------:|:---------:|:-----------|:----------------|:------------------------:|
| 2.5.1. | * | Tax table entry (TaxTableEntry) | | N/A |
| 2.5.1.1. | * | Tax type code (TaxType) | In this field, the type of tax must be indicated. Must be filled in with: "VAT" - Value Added Tax; "IS" - Stamp Duty; "NS" - Not subject to VAT or IS. | Text 3 |
| 2.5.1.2. | * | Tax country/region (TaxCountryRegion) | Must be filled in according to ISO 3166-1-alpha-2 standard. In the case of the Autonomous Regions of Madeira and Azores, it must be filled in with: "PT-AC" - Fiscal space of the Autonomous Region of the Azores; and "PT-MA" - Fiscal space of the Autonomous Region of Madeira. | Text 5 |
| 2.5.1.3. | * | Tax code (TaxCode) | In case field 2.5.1.1 - Tax type code (TaxType) = VAT, must be filled in with: "RED" - Reduced rate; "INT" - Intermediate rate; "NOR" - Normal rate; "ISE" - Exempt; and "OUT" - Others, applicable for special VAT regimes. In case field 2.5.1.1 - Tax type code (TaxType) = IS, must be filled in with: - The respective item code; - "ISE" - Exempt. In case of non-subjection, must be filled in with "NS". In receipts issued without itemized tax, must be filled in with "NA". | Text 10 |
| 2.5.1.4. | * | Description (Description) | In the case of stamp duty, must be filled in with the description of the respective item. | Text 255 |
| 2.5.1.5. | | Tax expiration date (TaxExpirationDate) | Last legal date of application of the tax rate, in case of change of the same, during the taxation period. | Date |
| 2.5.1.6. | ** | Tax percentage (TaxPercentage) | Filling is mandatory in case it is a tax percentage. In case of exemption or non-subjection to tax, use this field with the value "0" (zero). | Decimal |
| 2.5.1.7. | ** | Tax amount (TaxAmount) | Filling is mandatory in case it is a fixed stamp duty amount. | Monetary |

#### 3. Accounting movements (GeneralLedgerEntries)

In this table are recorded the accounting movements corresponding to the export period to which the SAF-T (PT) refers, and the opening movements should not be exported, as these are only reflected at the level of table 2.1. - Chart of accounts table (GeneralLedgerAccounts) in fields 2.1.2.3. - Opening debit balance (OpeningDebitBalance) and 2.1.2.4. - Opening credit balance (OpeningCreditBalance).

| Field Index | Mandatory | Field Name | Technical Notes | Format (validated in xsd) |
|:----------:|:---------:|:-----------|:----------------|:------------------------:|
| 3.1. | * | Number of accounting movement records (NumberOfEntries) | | Integer |
| 3.2. | * | Total debits (TotalDebit) | Sum of all debit movements of the selected period, recorded in field 3.4.3.11.1.6. Debit amount (DebitAmount). | Monetary |
| 3.3. | * | Total credits (TotalCredit) | Sum of all credit movements of the selected period, recorded in field 3.4.3.11.2.6. Credit amount (CreditAmount). | Monetary |
| 3.4. | | Journal (Journal) | | N/A |
| 3.4.1. | * | Journal ID (JournalID) | | Text 30 |
| 3.4.2. | * | Description (Description) | | Text 60 |
| 3.4.3. | | Transaction (Transaction) | | N/A |
| 3.4.3.1. | * | Unique accounting movement key (TransactionID) | Must be constructed in such a way as to be unique and to correspond to the accounting document number, which is used to detect the physical document in the archive, therefore, it must result from a concatenation, separated by spaces, between the following values: document date, journal identifier and document archive number (TransactionDate, JournalID and DocArchivalNumber). | Text 70 |
| 3.4.3.2. | * | Accounting period (Period) | The month number of the taxation period must be indicated, from "1" to "12", counted from its beginning. It can also be filled in with "13", "14", "15" or "16" for movements made in the last month of the taxation period, related to the determination of the result. Example: movements of inventory determinations, depreciations, adjustments or results determinations. | Integer |
| 3.4.3.3. | * | Transaction date (TransactionDate) | The date printed on the document that serves as support for the record must be indicated. | Date |
| 3.4.3.4. | * | User ID who recorded the movement (SourceID) | | Text 30 |
| 3.4.3.5. | * | Movement description (Description) | | Text 200 |
| 3.4.3.6. | * | Document archive number (DocArchivalNumber) | The document number within the journal must be indicated, which enables access to the originating document of the record. | Text 20 |
| 3.4.3.7. | * | Accounting movement type (TransactionType) | Must be filled in with: "N" - Normal; "R" - Regularizations of the taxation period; "A" - Results determination; "J" - Adjustment movements. | Text 1 |
| 3.4.3.8. | * | Accounting movement date (GLPostingDate) | This field must be filled in with the relevant date for accounting purposes. | Date |
| 3.4.3.9. | ** | Customer ID (CustomerID) | Filling is mandatory in case the customer is non-resident or the transaction constitutes a sale that must appear in annex O of the IES / annual declaration or that must appear in annex I of the periodic VAT declaration. The key of the record in table 2.2 - Customer table (Customer), contained in field 2.2.1 - Unique customer identifier (CustomerID) must be indicated. | Text 30 |
| 3.4.3.10. | ** | Supplier ID (SupplierID) | Filling is mandatory in case the supplier is non-resident or the transaction constitutes a purchase that must appear in annex P of the IES / annual declaration. The key of the record in table 2.3 - Supplier table (Supplier), contained in field 2.3.1 - Unique supplier identifier (SupplierID) must be indicated. | Text 30 |
| 3.4.3.11. | * | Lines (Lines) | Aggregating structure of movement lines, which must include at least one debit line and one credit line, without imposition in terms of ordering. | N/A |
| 3.4.3.11.1. | * | Debit line (DebitLine) | | N/A |
| 3.4.3.11.1.1. | * | Line record ID (RecordID) | The unique key of that line's record in the document must be indicated. | Text 30 |
| 3.4.3.11.1.2. | * | Account code (AccountID) | | Text 30 |
| 3.4.3.11.1.3. | | Source document ID (SourceDocumentID) | The type and number of the commercial document related to this line must be indicated. In case the application is integrated, the numbering structure of the source field must be used. | Text 60 |
| 3.4.3.11.1.4. | * | System entry date (SystemEntryDate) | Movement record to the second. Date and time type: "YYYY-MM-DDThh:mm:ss". | Date and Time |
| 3.4.3.11.1.5. | * | Document line description (Description) | | Text 200 |
| 3.4.3.11.1.6. | * | Debit amount (DebitAmount) | | Monetary |
| 3.4.3.11.2. | * | Credit line (CreditLine) | | N/A |
| 3.4.3.11.2.1. | * | Line record ID (RecordID) | The unique key of that line's record in the document must be indicated. | Text 30 |
| 3.4.3.11.2.2. | * | Account code (AccountID) | | Text 30 |
| 3.4.3.11.2.3. | | Source document ID (SourceDocumentID) | The type and number of the commercial document related to this line must be indicated. In case the application is integrated, the numbering structure of the source field must be used. | Text 60 |
| 3.4.3.11.2.4. | * | System entry date (SystemEntryDate) | Movement record to the second. Date and time type: "YYYY-MM-DDThh:mm:ss". | Date and Time |
| 3.4.3.11.2.5. | * | Document line description (Description) | | Text 200 |
| 3.4.3.11.2.6. | * | Credit amount (CreditAmount) | | Monetary |

#### 4. Commercial documents (SourceDocuments)

In these tables, lines without fiscal relevance should not be exported, namely technical descriptions, assembly instructions and warranty conditions.

The internal code of the document type cannot be used in different types of documents, regardless of the table in which it should be exported.

#### 4.1. Commercial documents to customers (SalesInvoices)

All sales documents and rectifying documents issued by the company must be included in this table, including canceled documents, duly marked, to allow verification of the sequentiality of numbering of documents within each document series, which must have at least annual numbering.

| Field Index | Mandatory | Field Name | Technical Notes | Format (validated in xsd) |
|:----------:|:---------:|:-----------|:----------------|:------------------------:|
| 4.1. | | Commercial documents to customers (SalesInvoices) | The documents indicated in field 4.1.4.8. - Document type (InvoiceType) must be exported. | N/A |
| 4.1.1. | * | Number of commercial document records (NumberOfEntries) | Must contain the total number of documents, including documents whose value in field 4.1.4.3.1. - Current status (InvoiceStatus) is of type "A" or "F". | Integer |
| 4.1.2. | * | Total debits (TotalDebit) | Must contain the control sum of field 4.1.4.19.13. - Debit amount (DebitAmount), excluding documents where field 4.1.4.3.1. - Current document status (InvoiceStatus) is of type "A" or "F". | Monetary |
| 4.1.3. | * | Total credits (TotalCredit) | Must contain the control sum of field 4.1.4.19.14. - Credit amount (CreditAmount), excluding documents where field 4.1.4.3.1. - Current document status (InvoiceStatus) is of type "A" or "F". | Monetary |
| 4.1.4. | | Sales document (Invoice) | | N/A |
| 4.1.4.1. | * | Unique sales document identification (InvoiceNo) | This identification is composed sequentially of the following elements: the internal code of the document type assigned by the application, a space, the document series identifier, a forward slash (/) and the sequential number of that document within that series. No records with the same identification can exist. The same internal document type code cannot be used in different types of documents. | Text 60 |
| 4.1.4.2. | * | Unique Document Code (ATCUD) | This field must contain the Unique Document Code. The field must be filled in with "0" (zero) until its regulation. | Text 100 |
| 4.1.4.3. | * | Document status (DocumentStatus) | | N/A |
| 4.1.4.3.1. | * | Current document status (InvoiceStatus) | Must be filled in with: "N" - Normal; "S" - Self-billing; "A" - Canceled document; "R" - Document summarizing other documents created in other applications and generated in this application; "F" - Invoiced document. | Text 1 |
| 4.1.4.3.2. | * | Document current status date and time (InvoiceStatusDate) | Date of the last status recording to the second. Date and time type: "YYYY-MM-DDThh:mm:ss". | Date and Time |
| 4.1.4.3.3. | | Reason for status change (Reason) | The reason that led to the change in document status must be indicated. | Text 50 |
| 4.1.4.3.4. | * | User code (SourceID) | User responsible for the current status of the document. | Text 30 |
| 4.1.4.3.5. | * | Document source (SourceBilling) | Must be filled in with: "P" - Document produced in the application; "I" - Document integrated and produced in another application; "M" - Document from recovery or manual issuance. | Text 1 |
| 4.1.4.4. | * | Document hash (Hash) | Signature under the terms of Ordinance No. 363/2010, of June 23. The field must be filled in with "0" (zero) if there is no certification requirement. | Text 172 |
| 4.1.4.5. | * | Hash control (HashControl) | Version of the private key used in creating the signature of field 4.1.4.4. - Document hash (Hash). The field must be filled in with "0" (zero) if the document is generated by an uncertified program. | Text 70 |
| 4.1.4.6. | | Accounting period (Period) | The month of the taxation period from "1" to "12" must be indicated, counted from the date of its beginning. | Integer |
| 4.1.4.7. | * | Invoice date (InvoiceDate) | Sales document issue date. | Date |
| 4.1.4.8. | * | Document type (InvoiceType) | Must be filled in with: "FT" - Invoice, issued under the terms of article 36 of the VAT Code; "FS" - Simplified invoice, issued under the terms of article 40 of the VAT Code; "FR" - Invoice-receipt; "ND" - Debit note; "NC" - Credit note; "VD" - Cash sale and invoice/receipt; (a) "TV" - Sales slip; (a) "TD" - Return slip; (a) "AA" - Asset disposal; (a) "DA" - Asset return. (a) For the Insurance sector when it should not appear in table 4.3. Goods or service provision checking documents (WorkingDocuments), it can also be filled in with: "RP" - Premium or premium receipt; "RE" - Cancellation or cancellation receipt; "CS" - Co-insurance allocation; "LD" - Lead co-insurer allocation; "RA" - Accepted reinsurance. (a) For data up to 2012-12-31. | Text 2 |
| 4.1.4.9. | * | Special regimes (SpecialRegimes) | | N/A |
| 4.1.4.9.1. | * | Self-billing indicator (SelfBillingIndicator) | Must be filled in with "1" if it concerns self-billing and with "0" (zero) otherwise. | Integer |
| 4.1.4.9.2. | * | Cash VAT scheme indicator (CashVATSchemeIndicator) | Indicator of the existence of adherence to the Cash VAT scheme. Must be filled in with "1" if there is adherence and with "0" (zero) otherwise. | Integer |
| 4.1.4.9.3. | * | Third parties billing indicator (ThirdPartiesBillingIndicator) | Must be filled in with "1" if it concerns billing issued in the name and on behalf of third parties and with "0" (zero) otherwise. | Integer |
| 4.1.4.10. | * | User code (SourceID) | User who generated the document. | Text 30 |
| 4.1.4.11. | | CAE code (EACCode) | The CAE code of the activity related to the issuance of this document must be indicated. | Text 5 |
| 4.1.4.12. | * | Document recording date (SystemEntryDate) | Date of record recording to the second, at the time of signing. Date and time type: "YYYY-MM-DDThh:mm:ss". | Date and Time |
| 4.1.4.13. | ** | Transaction identifier (TransactionID) | Filling is mandatory in the case of an integrated accounting and billing system, even if the file type (TaxAccountingBasis) should not contain the accounting tables. The unique key from table 3. Accounting movements (GeneralLedgerEntries) of the transaction where this document was posted must be indicated, respecting the rule defined there for field 3.4.3.1 - Unique accounting movement key (TransactionID). | Text 70 |
| 4.1.4.14. | * | Customer identifier (CustomerID) | Unique key from table 2.2. - Customer table (Customer) respecting the rule defined there for field 2.2.1. - Unique customer identifier (CustomerID). | Text 30 |
| 4.1.4.15. | | Delivery location (ShipTo) | Information about the location and date of delivery where the sold items are made available to the customer, or to whom they indicate in the case of triangular operations. | N/A |
| 4.1.4.15.1. | | Delivery identifier (DeliveryID) | In the case of triangular operations, the delivery location and TIN of who received must be indicated. | Text 255 |
| 4.1.4.15.2. | | Delivery date (DeliveryDate) | In the insurance sector, this field must be filled in with the end date of the risk coverage period. | Date |
| 4.1.4.15.3. | | Destination warehouse identifier (WarehouseID) | | Text 50 |
| 4.1.4.15.4. | | Goods location in destination warehouse (LocationID) | | Text 30 |
| 4.1.4.15.5. | | Address (Address) | | N/A |
| 4.1.4.15.5.1. | | Building number (BuildingNumber) | | Text 10 |
| 4.1.4.15.5.2. | | Street name (StreetName) | | Text 200 |
| 4.1.4.15.5.3. | * | Address detail (AddressDetail) | Must include street name, building number and floor, if applicable. | Text 210 |
| 4.1.4.15.5.4. | * | City (City) | | Text 50 |
| 4.1.4.15.5.5. | * | Postal code (PostalCode) | | Text 20 |
| 4.1.4.15.5.6. | | Region (Region) | | Text 50 |
| 4.1.4.15.5.7. | * | Country (Country) | Must be filled in according to ISO 3166-1-alpha-2 standard. | Text 2 |
| 4.1.4.16. | | Loading location (ShipFrom) | Information about the location and date of loading where the shipment of sold items to the customer begins. | N/A |
| 4.1.4.16.1. | | Delivery identifier (DeliveryID) | The registration number of the transport vehicle or the shipping method used must be indicated, e.g., express mail, etc. | Text 255 |
| 4.1.4.16.2. | | Reception date (DeliveryDate) | In the insurance sector, this field must be filled in with the start date of the risk coverage period. | Date |
| 4.1.4.16.3. | | Source warehouse identifier (WarehouseID) | | Text 50 |
| 4.1.4.16.4. | | Goods location in source warehouse (LocationID) | | Text 30 |
| 4.1.4.16.5. | | Address (Address) | | N/A |
| 4.1.4.16.5.1. | | Building number (BuildingNumber) | | Text 10 |
| 4.1.4.16.5.2. | | Street name (StreetName) | | Text 200 |
| 4.1.4.16.5.3. | * | Address detail (AddressDetail) | Must include street name, building number and floor, if applicable. | Text 210 |
| 4.1.4.16.5.4. | * | City (City) | | Text 50 |
| 4.1.4.16.5.5. | * | Postal code (PostalCode) | | Text 20 |
| 4.1.4.16.5.6. | | Region (Region) | | Text 50 |
| 4.1.4.16.5.7. | * | Country (Country) | Must be filled in according to ISO 3166-1-alpha-2 standard. | Text 2 |
| 4.1.4.17. | | Transport end date and time (MovementEndTime) | Date and time type: "YYYY-MM-DDThh:mm:ss" where "ss" can be "00" in the absence of concrete information. | Date and time |
| 4.1.4.18. | ** | Transport start date and time (MovementStartTime) | Date and time type: "YYYY-MM-DDThh:mm:ss" where "ss" can be "00" in the absence of concrete information. Filling is mandatory if the document also serves as a transport document according to the provisions of the Goods in Circulation Regime, approved by Decree-Law No. 147/2003, of July 11. | Date and time |
| 4.1.4.19. | * | Line (Line) | | N/A |
| 4.1.4.19.1. | * | Line number (LineNumber) | Lines must be exported in the same order as they appear in the original document. | Integer |
| 4.1.4.19.2. | | Source document reference (OrderReferences) | If there is a need to make more than one reference, this structure can be generated as many times as necessary. | N/A |
| 4.1.4.19.2.1. | | Source document number (OriginatingON) | The type, series and number of the document that supports this issuance must be indicated. If the document is contained in the SAF-T (PT), the numbering structure of the source field must be used. | Text 60 |
| 4.1.4.19.2.2. | | Source document date (OrderDate) | | Date |
| 4.1.4.19.3. | * | Product or service identifier (ProductCode) | Key of the record in table 2.4. Products/services table (Product), in field 2.4.2. Product or service identifier (ProductCode). | Text 60 |
| 4.1.4.19.4. | * | Product or service description (ProductDescription) | Description of the invoice line, linked to table 2.4. - Products/services table (Product), in field 2.4.4. Product or service description (ProductDescription). | Text 200 |
| 4.1.4.19.5. | * | Quantity (Quantity) | | Decimal |
| 4.1.4.19.6. | * | Unit of measure (UnitOfMeasure) | | Text 20 |
| 4.1.4.19.7. | * | Unit price (UnitPrice) | Unit price deducted from line and header discounts, without including taxes. Must be filled in with "0.00" if there is an obligation to fill in field 4.1.4.19.8. Unit tax base (TaxBase). | Monetary |
| 4.1.4.19.8. | ** | Unit tax base (TaxBase) | Unit tax base that does not contribute to the Total document without taxes (NetTotal). This value is what serves as the basis for calculating the line taxes. The sign (debit or credit) with which the tax thus calculated contributes to TaxPayable results from the existence in the line of DebitAmount or CreditAmount. | Monetary |
| 4.1.4.19.9. | * | Goods shipping or service provision date (TaxPointDate) | Date of goods shipping or service provision. Must be filled in with the associated delivery note date, if it exists. If there is more than one delivery note, the date of the oldest one must be indicated. | Date |
| 4.1.4.19.10. | | Invoice references (References) | References to invoices in their rectifying documents. If there is a need to make more than one reference, this structure can be generated as many times as necessary. | N/A |
| 4.1.4.19.10.1. | | Reference (Reference) | Reference to the invoice or simplified invoice, through its unique identification in the systems where it exists. The numbering structure of the source field must be used. | Text 60 |
| 4.1.4.19.10.2. | | Reason (Reason) | Must be filled in with the reason for issuance. | Text 50 |
| 4.1.4.19.11. | * | Line description (Description) | Description of the document line. | Text 200 |
| 4.1.4.19.12. | | Product serial number (ProductSerialNumber) | | N/A |
| 4.1.4.19.12.1. | * | Serial number (SerialNumber) | Must include the product serial number that appears in the document. Ex.: VIN, IMEI, ISSN, ISAN. If there is a need to make more than one reference, this field can be generated as many times as necessary. | Text 100 |
| 4.1.4.19.13. | ** | Debit amount (DebitAmount) | Line value, without tax, of documents to be posted to debit in the sales account. This value is deducted from line and header discounts. | Monetary |
| 4.1.4.19.14. | ** | Credit amount (CreditAmount) | Line value, without tax, of documents to be posted to credit in the sales account. This value is deducted from line and header discounts. | Monetary |
| 4.1.4.19.15. | * | Tax rate (Tax) | | N/A |
| 4.1.4.19.15.1. | * | Tax type code (TaxType) | In this field, the type of tax must be indicated. Must be filled in with: "VAT" - Value Added Tax; "IS" - Stamp Duty; "NS" - Not subject to VAT or IS. | Text 3 |
| 4.1.4.19.15.2. | * | Tax country/region (TaxCountryRegion) | Must be filled in according to ISO 3166-1-alpha-2 standard. In the case of the Autonomous Regions of Madeira and Azores, must be filled in with: "PT-AC" - Fiscal space of the Autonomous Region of the Azores; "PT-MA" - Fiscal space of the Autonomous Region of Madeira. | Text 5 |
| 4.1.4.19.15.3. | * | Tax code (TaxCode) | In case field 4.1.4.19.15.1. Tax type code (TaxType) = VAT, must be filled in with: "RED" - Reduced rate; "INT" - Intermediate rate; "NOR" - Normal rate; "ISE" - Exempt; "OUT" - Others, applicable for special VAT regimes. In case field 4.1.4.19.15.1. Tax type code (TaxType) = IS, must be filled in with: - The respective item code; - "ISE" - Exempt. In case of non-subjection, must be filled in with "NS". | Text 10 |
| 4.1.4.19.15.4. | ** | Tax rate percentage (TaxPercentage) | Filling is mandatory in case it is a tax percentage. The rate percentage corresponds to the tax applicable to field 4.1.4.19.13. - Debit amount (DebitAmount) or field 4.1.4.19.14. - Credit amount (CreditAmount). In case of exemption or non-subjection to tax, use this field with the value "0" (zero). | Decimal |
| 4.1.4.19.15.5. | ** | Tax amount (TaxAmount) | Filling is mandatory in case it is a fixed stamp duty unit amount. This value, multiplied by the quantity (Quantity) contributes to the tax payable amount (TaxPayable). | Monetary |
| 4.1.4.19.16. | ** | Tax exemption reason (TaxExemptionReason) | Filling is mandatory when fields 4.1.4.19.15.4. Tax rate percentage (TaxPercentage) or 4.1.4.19.15.5. Tax amount (TaxAmount) are equal to zero. The applicable legal provision must be referenced. This field must also be filled in in cases of non-subjection to the taxes referred to in table 2.5. - Tax table (TaxTable). | Text 60 |
| 4.1.4.19.17. | ** | Tax exemption code (TaxExemptionCode) | Must be filled in with the exemption or non-assessment reason code, which appears in the Software Integration Manual - Communication of Invoices to AT. Filling is mandatory when fields 4.1.4.19.15.4. Tax rate percentage (TaxPercentage) or 4.1.4.19.15.5. Tax amount (TaxAmount) are equal to zero. This field must also be filled in in cases of non-subjection to the taxes referred to in table 2.5. - Tax table (TaxTable). | Text 3 |
| 4.1.4.19.18. | | Line settlement amount (SettlementAmount) | Must reflect all discounts granted (the proportion of global discounts for this line and those specific to the same line) that affect the value of field 4.1.4.20.3. Total document with taxes (GrossTotal). | Monetary |
| 4.1.4.19.19. | | Customs information (CustomsInformation) | | N/A |
| 4.1.4.19.19.1. | | Administrative reference code (ARCNo) | Fill in with the code assigned after validation of the electronic administrative document e-DA. If there is a need to make more than one reference, this field can be generated as many times as necessary. | Text 21 |
| 4.1.4.19.19.2. | | Line special consumption tax amount (IECAmount) | Amount of special consumption tax contained in the taxable base of the document line if this is not discriminated in an independent line in the document, with the "Product or service indicator (ProductType)" = E. | Monetary |
| 4.1.4.20. | * | Document totals (DocumentTotals) | | N/A |
| 4.1.4.20.1. | * | Tax payable (TaxPayable) | | Monetary |
| 4.1.4.20.2. | * | Total document without taxes (NetTotal) | This field must not include the tax portions contained in table 2.5. - Tax table (TaxTable). | Monetary |
| 4.1.4.20.3. | * | Total document with taxes (GrossTotal) | This field must not reflect any withholding taxes contained in structure 4.1.4.21. Withholding tax (WithholdingTax). | Monetary |
| 4.1.4.20.4. | | Currency (Currency) | Must not be generated in case the document issuance is in euro. | N/A |
| 4.1.4.20.4.1. | * | Currency code (CurrencyCode) | In case of foreign currency must be filled in according to ISO 4217 standard. | Text 3 |
| 4.1.4.20.4.2. | * | Currency amount (CurrencyAmount) | Value of field 4.1.4.20.3. - Total document with taxes (GrossTotal) in the original currency of the document. | Monetary |
| 4.1.4.20.5. | | Settlement (Settlement) | Payment agreements or methods. If there is a need to make more than one reference, this structure can be generated as many times as necessary. | N/A |
| 4.1.4.20.5.1. | | Settlement discount (SettlementDiscount) | Must be filled in with future discount agreements to be applied on the present value. | Text 30 |
| 4.1.4.20.5.2. | | Settlement amount (SettlementAmount) | Represents the value agreed for future discount without affecting the present value of the document indicated in field 4.1.4.20.3. - Total document with taxes (GrossTotal). | Monetary |
| 4.1.4.20.5.3. | | Settlement date (SettlementDate) | The information to be included is the agreed date for payment with discount. | Date |
| 4.1.4.20.6. | | Payment (Payment) | In case of mixed payments, the amounts by type and payment date must be indicated. If there is a need to make more than one reference, this structure can be generated as many times as necessary. | N/A |
| 4.1.4.20.6.1. | | Payment mechanism (PaymentMechanism) | Must be filled in with: "CC" - Credit card; "CD" - Debit card; "CH" - Bank check; "CI" - International documentary credit; "CO" - Gift check or card; "CS" - Current account balance compensation; "DE" - Electronic money, for example resident in loyalty or points cards; "LC" - Commercial bill; "MB" - Payment references for Multibanco; "NU" - Cash; "OU" - Other means not indicated here; "PR" - Goods exchange; "TB" - Bank transfer or authorized direct debit; "TR" - Extra-salary compensation titles regardless of their support, for example, meal vouchers, education, etc. | Text 2 |
| 4.1.4.20.6.2. | * | Payment amount (PaymentAmount) | The amount by payment method must be indicated. | Monetary |
| 4.1.4.20.6.3. | * | Payment date (PaymentDate) | | Date |
| 4.1.4.21. | | Withholding tax (WithholdingTax) | If there is a need to make more than one reference, this structure can be generated as many times as necessary. | N/A |
| 4.1.4.21.1. | | Withholding tax type (WithholdingTaxType) | In this field, the type of tax withheld must be indicated, filling it with: "IRS" - Personal income tax; "IRC" - Corporate income tax; "IS" - Stamp duty. | Text 3 |
| 4.1.4.21.2. | | Withholding tax description (WithholdingTaxDescription) | The applicable legal provision must be indicated. In case of Tax type code (WithholdingTaxType) = IS, must be filled in with the respective item code. | Text 60 |
| 4.1.4.21.3. | * | Withholding tax amount (WithholdingTaxAmount) | The amount of tax withheld must be indicated. | Monetary |

#### 4.2. Goods movement documents (MovementOfGoods)

Documents must be exported, namely transport or delivery notes, that serve as transport documents, in accordance with the provisions of the Goods in Circulation Regime, approved by Decree-Law No. 147/2003, of July 11.

Those documents that must appear in table 4.1. Commercial documents to customers (SalesInvoices) and that also served as transport documents (for example invoices) should not be exported here.

| Field Index | Mandatory | Field Name | Technical Notes | Format (validated in xsd) |
|:----------:|:---------:|:-----------|:----------------|:------------------------:|
| 4.2. | | Goods movements (MovementOfGoods) | Documents must be exported, namely transport or delivery notes, that serve as transport documents, in accordance with the provisions of the Goods in Circulation Regime, approved by Decree-Law No. 147/2003, of July 11. Those documents that must appear in table 4.1. Commercial documents to customers (SalesInvoices) and that also served as transport documents (for example invoices) should not be exported here. | N/A |
| 4.2.1. | * | Number of goods movement line records (NumberOfMovementLines) | Must contain the total number of lines with fiscal relevance of the documents of the period made available, including the lines of documents where field 4.2.3.3.1. - Current document status (MovementStatus), is of type "A". | Integer |
| 4.2.2. | * | Total quantities moved (TotalQuantityIssued) | Must contain the control sum of field 4.2.3.21.5. - Quantity (Quantity) excluding the lines of documents where field 4.2.3.3.1. - Current document status (MovementStatus), is of type "A". | Decimal |
| 4.2.3. | | Goods movement document (StockMovement) | | N/A |
| 4.2.3.1. | * | Unique document identification (DocumentNumber) | This identification is composed sequentially of the following elements: the internal code of the document type assigned by the application, a space, the document series identifier, a forward slash (/) and the sequential number of that document within that series. No records with the same identification can exist in this field. The same internal document type code cannot be used in different types of documents. | Text 60 |
| 4.2.3.2. | * | Unique Document Code (ATCUD) | This field must contain the Unique Document Code. The field must be filled in with "0" (zero) until its regulation. | Text 100 |
| 4.2.3.3. | * | Document status (DocumentStatus) | | N/A |
| 4.2.3.3.1. | * | Current document status (MovementStatus) | Must be filled in with: "N" - Normal; "T" - On behalf of third parties; "A" - Canceled document; "F" - Invoiced document, even if partially, when for this document there is also in table 4.1. - Commercial documents to customers (SalesInvoices) the corresponding invoice or simplified invoice type. "R" - Document summarizing other documents created in other applications and generated in this application. | Text 1 |
| 4.2.3.3.2. | * | Document current status date and time (MovementStatusDate) | Date of the last status recording to the second. Date and time type: "YYYY-MM-DDThh:mm:ss". | Date and Time |
| 4.2.3.3.3. | | Reason for status change (Reason) | The reason that led to the change in document status must be indicated. | Text 50 |
| 4.2.3.3.4. | * | User code (SourceID) | User responsible for the current status of the document. | Text 30 |
| 4.2.3.3.5. | * | Document source (SourceBilling) | Must be filled in with: "P" - Document produced in the application; "I" - Document integrated and produced in another application; "M" - Document from recovery or manual issuance. | Text 1 |
| 4.2.3.4. | * | Document hash (Hash) | Signature under the terms of Ordinance No. 363/2010, of June 23. The field must be filled in with "0" (zero) if there is no certification requirement. | Text 172 |
| 4.2.3.5. | * | Hash control (HashControl) | Version of the private key used in creating the signature of field 4.2.3.4. - Document hash (Hash). The field must be filled in with "0" (zero) if the document is generated by an uncertified program. | Text 70 |
| 4.2.3.6. | | Accounting period (Period) | The month of the taxation period from "1" to "12" must be indicated, counted from the date of its beginning. | Integer |
| 4.2.3.7. | * | Goods movement date (MovementDate) | Transport document issue date. | Date |
| 4.2.3.8. | * | Movement type (MovementType) | Must be filled in with: "GR" - Delivery note; "GT" - Transport note (include here global notes); "GA" - Own fixed assets movement note; "GC" - Consignment note; "GD" - Return note or return slip. | Text 2 |
| 4.2.3.9. | * | Document recording date (SystemEntryDate) | Date of record recording to the second. Date and time type: "YYYY-MM-DDThh:mm:ss". | Date and Time |
| 4.2.3.10. | ** | Transaction identifier (TransactionID) | If it gives rise to an accounting record, filling is mandatory in the case of an integrated accounting and billing system, even if the file type (TaxAccountingBasis) should not contain the accounting tables. The unique key from table 3 - Accounting movements (GeneralLedgerEntries) where this document was posted must be indicated, respecting the rule defined there for field 3.4.3.1 - Unique accounting movement key (TransactionID). | Text 70 |
| 4.2.3.11. | ** | Customer identifier (CustomerID) | Unique key from table 2.2 - Customer table (Customer) respecting the rule defined there for field 2.2.1 - Unique customer identifier (CustomerID). In the case of notes where the recipient is not known, the generic customer provided in table 2.2. - Customer table (Customer) must be used. This field must also be filled in in the case of notes that title the transfer of goods from the sender themselves. | Text 30 |
| 4.2.3.12. | ** | Supplier identifier (SupplierID) | Unique key from table 2.3 - Supplier table (Supplier) respecting the rule defined there for field 2.3.1 - Unique supplier identifier (SupplierID), in the case of return notes or transport note of movable goods produced or assembled to order with materials that the work owner has supplied for that purpose (custom work). | Text 30 |
| 4.2.3.13. | * | User code (SourceID) | User who generated the document. | Text 30 |
| 4.2.3.14. | | CAE code (EACCode) | The CAE code of the activity related to the issuance of this document must be indicated. | Text 5 |
| 4.2.3.15. | | Document issuance reason (MovementComments) | | Text 60 |
| 4.2.3.16. | | Delivery location (ShipTo) | Information about the location and date of delivery where the goods are made available to the customer, or to whom they indicate in the case of triangular operations. | N/A |
| 4.2.3.16.1. | | Delivery identifier (DeliveryID) | In the case of triangular operations, the delivery location and TIN of who received must be indicated. | Text 255 |
| 4.2.3.16.2. | | Delivery date (DeliveryDate) | | Date |
| 4.2.3.16.3. | | Destination warehouse identifier (WarehouseID) | | Text 50 |
| 4.2.3.16.4. | | Goods location in destination warehouse (LocationID) | | Text 30 |
| 4.2.3.16.5. | | Address (Address) | | N/A |
| 4.2.3.16.5.1. | | Building number (BuildingNumber) | | Text 10 |
| 4.2.3.16.5.2. | | Street name (StreetName) | | Text 200 |
| 4.2.3.16.5.3. | * | Address detail (AddressDetail) | Must include street name, building number and floor, if applicable. | Text 210 |
| 4.2.3.16.5.4. | * | City (City) | | Text 50 |
| 4.2.3.16.5.5. | * | Postal code (PostalCode) | | Text 20 |
| 4.2.3.16.5.6. | | Region (Region) | | Text 50 |
| 4.2.3.16.5.7. | * | Country (Country) | Must be filled in according to ISO 3166-1-alpha-2 standard. | Text 2 |
| 4.2.3.17. | | Loading location (ShipFrom) | Information about the location and date of loading where the shipment of goods to the customer begins. | N/A |
| 4.2.3.17.1. | | Delivery identifier (DeliveryID) | The registration number of the transport vehicle or the shipping method used must be indicated, e.g., express mail, etc. | Text 255 |
| 4.2.3.17.2. | | Shipping date (DeliveryDate) | | Date |
| 4.2.3.17.3. | | Source warehouse identifier (WarehouseID) | | Text 50 |
| 4.2.3.17.4. | | Goods location in source warehouse (LocationID) | | Text 30 |
| 4.2.3.17.5. | | Address (Address) | | N/A |
| 4.2.3.17.5.1. | | Building number (BuildingNumber) | | Text 10 |
| 4.2.3.17.5.2. | | Street name (StreetName) | | Text 200 |
| 4.2.3.17.5.3. | * | Address detail (AddressDetail) | Must include street name, building number and floor, if applicable. | Text 210 |
| 4.2.3.17.5.4. | * | City (City) | | Text 50 |
| 4.2.3.17.5.5. | * | Postal code (PostalCode) | | Text 20 |
| 4.2.3.17.5.6. | | Region (Region) | | Text 50 |
| 4.2.3.17.5.7. | * | Country (Country) | Must be filled in according to ISO 3166-1-alpha-2 standard. | Text 2 |
| 4.2.3.18. | | Transport end date and time (MovementEndTime) | Date and time type: "YYYY-MM-DDThh:mm:ss" where "ss" can be "00" in the absence of concrete information. | Date and time |
| 4.2.3.19. | * | Transport start date and time (MovementStartTime) | Date and time type: "YYYY-MM-DDThh:mm:ss" where "ss" can be "00" in the absence of concrete information. | Date and time |
| 4.2.3.20. | | Document identification code (ATDocCodeID) | Document identification code assigned by AT, under the terms of DL 147/2003, of July 11. | Text 200 |
| 4.2.3.21. | * | Line (Line) | | N/A |
| 4.2.3.21.1. | * | Line number (LineNumber) | Lines must be exported in the same order as they appear in the original document. | Integer |
| 4.2.3.21.2. | | Source document reference (OrderReferences) | If there is a need to make more than one reference, this structure can be generated as many times as necessary. | N/A |
| 4.2.3.21.2.1. | | Source document number (OriginatingON) | If the document is contained in the SAF-T (PT), the numbering structure of the source field must be used. | Text 60 |
| 4.2.3.21.2.2. | | Source document date (OrderDate) | | Date |
| 4.2.3.21.3. | * | Product or service identifier (ProductCode) | Key of the record in table 2.4. Products/services table (Product), in field 2.4.2. Product or service identifier (ProductCode). | Text 60 |
| 4.2.3.21.4. | * | Product or service description (ProductDescription) | Description of the invoice line, linked to table 2.4. - Products/services table (Product), in field 2.4.4. Product or service description (ProductDescription). | Text 200 |
| 4.2.3.21.5. | * | Quantity (Quantity) | | Decimal |
| 4.2.3.21.6. | * | Unit of measure (UnitOfMeasure) | | Text 20 |
| 4.2.3.21.7. | * | Unit price (UnitPrice) | Unit price deducted from line and header discounts, without including taxes. When not valued in the database, must be filled in with "0.00". | Monetary |
| 4.2.3.21.8. | * | Line description (Description) | Description of the document line. | Text 200 |
| 4.2.3.21.9. | | Product serial number (ProductSerialNumber) | | N/A |
| 4.2.3.21.9.1. | * | Serial number (SerialNumber) | Must include the product serial number that appears in the document. Ex.: VIN, IMEI, ISSN, ISAN. If there is a need to make more than one reference, this field can be generated as many times as necessary. | Text 100 |
| 4.2.3.21.10. | ** | Debit amount (DebitAmount) | Field to be filled in in case of valued goods entry. When not valued in the database, must be filled in with "0.00". | Monetary |
| 4.2.3.21.11. | ** | Credit amount (CreditAmount) | Field to be filled in in case of valued goods exit. When not valued in the database, must be filled in with "0.00". | Monetary |
| 4.2.3.21.12. | ** | Tax rate (Tax) | This structure should only be created in documents valued in the database. | N/A |
| 4.2.3.21.12.1. | * | Tax type code (TaxType) | In this field, the type of tax must be indicated. Must be filled in with: "VAT" - Value Added Tax; "NS" - Not subject to VAT. | Text 3 |
| 4.2.3.21.12.2. | * | Tax country/region (TaxCountryRegion) | Must be filled in according to ISO 3166-1-alpha-2 standard. In the case of the Autonomous Regions of Madeira and Azores, must be filled in with: "PT-AC" - Fiscal space of the Autonomous Region of the Azores; "PT-MA" - Fiscal space of the Autonomous Region of Madeira. | Text 5 |
| 4.2.3.21.12.3. | * | Tax code (TaxCode) | Tax code in the tax table. Must be filled in with: "RED" - Reduced rate; "INT" - Intermediate rate; "NOR" - Normal rate; "ISE" - Exempt; "OUT" - Others, applicable for special VAT regimes. In case of non-subjection must be filled in with "NS". | Text 10 |
| 4.2.3.21.12.4. | * | Tax rate percentage (TaxPercentage) | Tax rate percentage corresponding to the tax applicable to field 4.2.3.21.10. Debit amount (DebitAmount) or field 4.2.3.21.11. - Credit amount (CreditAmount). In case of exemption or non-subjection to tax, use this field with the value "0" (zero). | Decimal |
| 4.2.3.21.13. | ** | Tax exemption reason (TaxExemptionReason) | Filling is mandatory when field 4.2.3.21.12.4. - Tax rate percentage (TaxPercentage) is equal to "0" (zero). The applicable legal provision must be referenced. This field must also be filled in in cases of non-subjection to the taxes referred to in table 2.5. - Tax table (TaxTable). | Text 60 |
| 4.2.3.21.14. | ** | Tax exemption code (TaxExemptionCode) | Must be filled in with the exemption or non-assessment reason code, which appears in the Software Integration Manual - Communication of Invoices to AT. Filling is mandatory when field 4.2.3.21.12.4. Tax rate percentage (TaxPercentage) is equal to zero. This field must also be filled in in cases of non-subjection to the taxes referred to in table 2.5. - Tax table (TaxTable). | Text 3 |
| 4.2.3.21.15. | | Line settlement amount (SettlementAmount) | Must reflect all discounts granted (the proportion of global discounts for this line and those specific to the same line) that affect the value of field 4.2.3.22.3. Total document with taxes (GrossTotal). | Monetary |
| 4.2.3.21.16. | | Customs information (CustomsInformation) | | N/A |
| 4.2.3.21.16.1. | | Administrative reference code (ARCNo) | Fill in with the code assigned after validation of the electronic administrative document e-DA. If there is a need to make more than one reference, this field can be generated as many times as necessary. | Text 21 |
| 4.2.3.21.16.2. | | Line special consumption tax amount (IECAmount) | Amount of special consumption tax contained in the taxable base of the document line if this is not discriminated in an independent line in the document, with the "Product or service indicator (ProductType)" = E. | Monetary |
| 4.2.3.22. | * | Document totals (DocumentTotals) | | N/A |
| 4.2.3.22.1. | * | Tax payable (TaxPayable) | When not valued in the database, must be filled in with "0.00". | Monetary |
| 4.2.3.22.2. | * | Total document without taxes (NetTotal) | This field must not include the tax portions contained in table 2.5. - Tax table (TaxTable). When not valued in the database, must be filled in with "0.00". | Monetary |
| 4.2.3.22.3. | * | Total document with taxes (GrossTotal) | When not valued in the database, must be filled in with "0.00". | Monetary |
| 4.2.3.22.4. | | Currency (Currency) | Must not be generated in case the document issuance is in euro. | N/A |
| 4.2.3.22.4.1. | * | Currency code (CurrencyCode) | In case of foreign currency must be filled in according to ISO 4217 standard. | Text 3 |
| 4.2.3.22.4.2. | * | Currency amount (CurrencyAmount) | Value of field 4.2.3.22.3. - Total document with taxes (GrossTotal) in the original currency of the document. | Monetary |
| 4.2.3.22.4.3. | * | Exchange rate (ExchangeRate) | The exchange rate used in the conversion to EUR must be indicated. | Decimal |

#### 4.3. Goods or service provision checking documents (WorkingDocuments)

In this table must be exported any other documents issued, regardless of their designation, that may be presented to the customer for checking goods or service provision, even if subject to subsequent billing.

Those documents that must appear in tables 4.1. Commercial documents to customers (SalesInvoices) or 4.2. - Goods movement documents (MovementOfGoods) should not be exported here.

| Field Index | Mandatory | Field Name | Technical Notes | Format (validated in xsd) |
|:----------:|:---------:|:-----------|:----------------|:------------------------:|
| 4.3. | | Checking documents (WorkingDocuments) | | N/A |
| 4.3.1. | * | Number of checking document records (NumberOfEntries) | Must contain the total number of documents, including documents whose value in field 4.3.4.3.1. Current status (WorkStatus) is of type "A". | Integer |
| 4.3.2. | * | Total debits (TotalDebit) | Must contain the control sum of field 4.3.4.14.13. - Debit amount (DebitAmount), excluding documents where field 4.3.4.3.1. - Current document status (WorkStatus) is of type "A". | Monetary |
| 4.3.3. | * | Total credits (TotalCredit) | Must contain the control sum of field 4.3.4.14.14. - Credit amount (CreditAmount), excluding documents where field 4.3.4.3.1. - Current document status (WorkStatus) is of type "A". | Monetary |
| 4.3.4. | | Checking document (WorkDocument) | | N/A |
| 4.3.4.1. | * | Unique document identification (DocumentNumber) | This identification is composed sequentially of the following elements: the internal code of the document type assigned by the application, a space, the document series identifier, a forward slash (/) and the sequential number of that document within that series. No records with the same identification can exist in this field. The same internal document type code cannot be used in different types of documents. | Text 60 |
| 4.3.4.2. | * | Unique Document Code (ATCUD) | This field must contain the Unique Document Code. The field must be filled in with "0" (zero) until its regulation. | Text 100 |
| 4.3.4.3. | * | Document status (DocumentStatus) | | N/A |
| 4.3.4.3.1. | * | Current document status (WorkStatus) | Must be filled in with: "N" - Normal; "A" - Canceled document; "F" - Invoiced document, even if partially, when for this document there is also in table 4.1. - Commercial documents to customers (SalesInvoices) the corresponding invoice or simplified invoice type. | Text 1 |
| 4.3.4.3.2. | * | Document current status date and time (WorkStatusDate) | Date of the last status recording to the second. Date and time type: "YYYY-MM-DDThh:mm:ss". | Date and Time |
| 4.3.4.3.3. | | Reason for status change (Reason) | The reason that led to the change in document status must be indicated. | Text 50 |
| 4.3.4.3.4. | * | User code (SourceID) | User responsible for the current status of the document. | Text 30 |
| 4.3.4.3.5. | * | Document source (SourceBilling) | Must be filled in with: "P" - Document produced in the application; "I" - Document integrated and produced in another application; "M" - Document from recovery or manual issuance. | Text 1 |
| 4.3.4.4. | * | Document hash (Hash) | Signature under the terms of Ordinance No. 363/2010, of June 23. The field must be filled in with "0" (zero) if there is no certification requirement. | Text 172 |
| 4.3.4.5. | * | Hash control (HashControl) | Version of the private key used in creating the signature of field 4.3.4.4. - Document hash (Hash). The field must be filled in with "0" (zero) if the document is generated by an uncertified program. | Text 70 |
| 4.3.4.6. | | Accounting period (Period) | The month of the taxation period from "1" to "12" must be indicated, counted from the date of its beginning. | Integer |
| 4.3.4.7. | * | Document date (WorkDate) | Operative document issue date. | Date |
| 4.3.4.8. | * | Document type (WorkType) | Must be filled in with: "CM" - Table consultations; "CC" - Consignment credit; "FC" - Consignment invoice under the terms of art. 38 of the VAT code; "FO" - Work sheets; "NE" - Order note; "OU" - Others; "OR" - Budgets; "PF" - Pro-forma; "DC" - Documents issued that may be presented to the customer for checking goods or service provision (for data up to 2017-06-30). For the Insurance sector when for the document types identified below there should also exist in table 4.1 Commercial documents to customers (SalesInvoices) the corresponding invoice or invoice rectifying document, it can also be filled in with: "RP" - Premium or premium receipt; "RE" - Cancellation or cancellation receipt; "CS" - Co-insurance allocation; "LD" - Lead co-insurer allocation; "RA" - Accepted reinsurance. | Text 2 |
| 4.3.4.9. | * | User code (SourceID) | User who generated the document. | Text 30 |
| 4.3.4.10. | | CAE code (EACCode) | The CAE code of the activity related to the issuance of this document must be indicated. | Text 5 |
| 4.3.4.11. | * | Document recording date (SystemEntryDate) | Date of record recording to the second. Date and time type: "YYYY-MM-DDThh:mm:ss". | Date and Time |
| 4.3.4.12. | ** | Transaction identifier (TransactionID) | If it gives rise to an accounting record, filling is mandatory in the case of an integrated accounting and billing system, even if the file type (TaxAccountingBasis) should not contain the accounting tables. The unique key from table 3. Accounting movements (GeneralLedgerEntries) where this document was posted must be indicated, respecting the rule defined there for field 3.4.3.1 - Unique accounting movement key (TransactionID). | Text 70 |
| 4.3.4.13. | * | Customer identifier (CustomerID) | Unique key from table 2.2. - Customer table (Customer) respecting the rule defined there for field 2.2.1. - Unique customer identifier (CustomerID). | Text 30 |
| 4.3.4.14. | | References to source documents (References) | References to source documents. | N/A |
| 4.3.4.14.1. | | Reference (Reference) | This structure can be generated as many times as necessary. | N/A |
| 4.3.4.14.1.1. | | Reference origin (Reference) | Must be filled in with: "P" - Document produced in the application; "I" - Document integrated and produced in another application; "M" - Document from recovery or manual issuance. | Text 1 |
| 4.3.4.14.1.2. | | Reference reason (Reason) | Reason for issuing this document. | Text 50 |
| 4.3.4.14.1.3. | | Referenced document number (CreditReference) | Document number being referenced. | Text 60 |
| 4.3.4.14.1.4. | | Referenced document date (CreditDate) | Date of the referenced document. | Date |
| 4.3.4.15. | | Shipping address (ShipTo) | Information about the delivery location and date of delivery of goods to the customer. | N/A |
| 4.3.4.15.1. | | Delivery identifier (DeliveryID) | Registration number of the transport vehicle or shipping method used. | Text 255 |
| 4.3.4.15.2. | | Delivery date (DeliveryDate) | Date of delivery of goods to the customer. | Date |
| 4.3.4.15.3. | | Warehouse identifier (WarehouseID) | Identifier of the warehouse where the goods are delivered. | Text 50 |
| 4.3.4.15.4. | | Location identifier (LocationID) | Location of goods in the warehouse. | Text 30 |
| 4.3.4.15.5. | | Address (Address) | Address details. | N/A |
| 4.3.4.15.5.1. | | Building number (BuildingNumber) | Building number. | Text 10 |
| 4.3.4.15.5.2. | | Street name (StreetName) | Street name. | Text 200 |
| 4.3.4.15.5.3. | | Address detail (AddressDetail) | Must include street name, building number and floor, if applicable. | Text 210 |
| 4.3.4.15.5.4. | | City (City) | City. | Text 50 |
| 4.3.4.15.5.5. | | Postal code (PostalCode) | Postal code. | Text 20 |
| 4.3.4.15.5.6. | | Region (Region) | Region. | Text 50 |
| 4.3.4.15.5.7. | | Country (Country) | Must be filled in according to ISO 3166-1-alpha-2 standard. | Text 2 |
| 4.3.4.16. | * | Line (Line) | | N/A |
| 4.3.4.16.1. | * | Line number (LineNumber) | Lines must be exported in the same order as they appear in the original document. | Integer |
| 4.3.4.16.2. | | Source document reference (OrderReferences) | If there is a need to make more than one reference, this structure can be generated as many times as necessary. | N/A |
| 4.3.4.16.2.1. | | Source document number (OriginatingON) | If the document is contained in the SAF-T (PT), the numbering structure of the source field must be used. | Text 60 |
| 4.3.4.16.2.2. | | Source document date (OrderDate) | | Date |
| 4.3.4.16.3. | * | Product or service identifier (ProductCode) | Key of the record in table 2.4. Products/services table (Product), in field 2.4.2. Product or service identifier (ProductCode). | Text 60 |
| 4.3.4.16.4. | * | Product or service description (ProductDescription) | Description of the invoice line, linked to table 2.4. - Products/services table (Product), in field 2.4.4. Product or service description (ProductDescription). | Text 200 |
| 4.3.4.16.5. | * | Quantity (Quantity) | | Decimal |
| 4.3.4.16.6. | * | Unit of measure (UnitOfMeasure) | | Text 20 |
| 4.3.4.16.7. | * | Unit price (UnitPrice) | Unit price deducted from line and header discounts, without including taxes. | Monetary |
| 4.3.4.16.8. | ** | Unit tax base (TaxBase) | Unit tax base that does not contribute to the Total document without taxes (NetTotal). This value is what serves as the basis for calculating the line taxes. | Monetary |
| 4.3.4.16.9. | * | Line description (Description) | Description of the document line. | Text 200 |
| 4.3.4.16.10. | | Product serial number (ProductSerialNumber) | | N/A |
| 4.3.4.16.10.1. | * | Serial number (SerialNumber) | Must include the product serial number that appears in the document. Ex.: VIN, IMEI, ISSN, ISAN. If there is a need to make more than one reference, this field can be generated as many times as necessary. | Text 100 |
| 4.3.4.16.11. | ** | Debit amount (DebitAmount) | Line value, without tax, of documents to be posted to debit. This value is deducted from line and header discounts. | Monetary |
| 4.3.4.16.12. | ** | Credit amount (CreditAmount) | Line value, without tax, of documents to be posted to credit. This value is deducted from line and header discounts. | Monetary |
| 4.3.4.16.13. | * | Tax rate (Tax) | | N/A |
| 4.3.4.16.13.1. | * | Tax type code (TaxType) | In this field, the type of tax must be indicated. Must be filled in with: "VAT" - Value Added Tax; "IS" - Stamp Duty; "NS" - Not subject to VAT or IS. | Text 3 |
| 4.3.4.16.13.2. | * | Tax country/region (TaxCountryRegion) | Must be filled in according to ISO 3166-1-alpha-2 standard. In the case of the Autonomous Regions of Madeira and Azores, must be filled in with: "PT-AC" - Fiscal space of the Autonomous Region of the Azores; "PT-MA" - Fiscal space of the Autonomous Region of Madeira. | Text 5 |
| 4.3.4.16.13.3. | * | Tax code (TaxCode) | Must be filled in with: "RED" - Reduced rate; "INT" - Intermediate rate; "NOR" - Normal rate; "ISE" - Exempt; "OUT" - Others, applicable for special VAT regimes. In case of non-subjection must be filled in with "NS". | Text 10 |
| 4.3.4.16.13.4. | ** | Tax rate percentage (TaxPercentage) | Tax rate percentage corresponding to the tax applicable to field 4.3.4.16.11. - Debit amount (DebitAmount) or field 4.3.4.16.12. - Credit amount (CreditAmount). In case of exemption or non-subjection to tax, use this field with the value "0" (zero). | Decimal |
| 4.3.4.16.13.5. | ** | Tax amount (TaxAmount) | Filling is mandatory in case it is a fixed stamp duty unit amount. | Monetary |
| 4.3.4.16.14. | ** | Tax exemption reason (TaxExemptionReason) | Filling is mandatory when fields 4.3.4.16.13.4. - Tax rate percentage (TaxPercentage) or 4.3.4.16.13.5. - Tax amount (TaxAmount) are equal to zero. The applicable legal provision must be referenced. | Text 60 |
| 4.3.4.16.15. | ** | Tax exemption code (TaxExemptionCode) | Must be filled in with the exemption or non-assessment reason code, which appears in the Software Integration Manual - Communication of Invoices to AT. | Text 3 |
| 4.3.4.16.16. | | Line settlement amount (SettlementAmount) | Must reflect all discounts granted (the proportion of global discounts for this line and those specific to the same line) that affect the value of field 4.3.4.17.3. Total document with taxes (GrossTotal). | Monetary |
| 4.3.4.17. | * | Document totals (DocumentTotals) | | N/A |
| 4.3.4.17.1. | * | Tax payable (TaxPayable) | | Monetary |
| 4.3.4.17.2. | * | Total document without taxes (NetTotal) | This field must not include the tax portions contained in table 2.5. - Tax table (TaxTable). | Monetary |
| 4.3.4.17.3. | * | Total document with taxes (GrossTotal) | | Monetary |
| 4.3.4.17.4. | | Currency (Currency) | Must not be generated in case the document issuance is in euro. | N/A |
| 4.3.4.17.4.1. | * | Currency code (CurrencyCode) | In case of foreign currency must be filled in according to ISO 4217 standard. | Text 3 |
| 4.3.4.17.4.2. | * | Currency amount (CurrencyAmount) | Value of field 4.3.4.17.3. - Total document with taxes (GrossTotal) in the original currency of the document. | Monetary |
| 4.3.4.17.4.3. | * | Exchange rate (ExchangeRate) | The exchange rate used in the conversion to EUR must be indicated. | Decimal |

#### 4.4. Payment receipts (Payments)

In this table must be exported the issued receipts, created after the entry into force of this structure.

| Field Index | Mandatory | Field Name | Technical Notes | Format (validated in xsd) |
|:----------:|:---------:|:-----------|:----------------|:------------------------:|
| 4.4. | | Payment receipts issued (Payments) | Must be exported the documents indicated in field 4.4.4.6. - Receipt type (PaymentType). | N/A |
| 4.4.1. | * | Number of issued receipt records (NumberOfEntries) | Must contain the total number of issued receipts, including documents whose value in field 4.4.4.9.1. - Current receipt status (PaymentStatus) is of type "A". | Integer |
| 4.4.2. | * | Total debits (TotalDebit) | Must contain the control sum of field 4.4.4.14.4. - Debit amount (DebitAmount), excluding documents where field 4.4.4.9.1. Current receipt status (PaymentStatus) is of type "A". | Monetary |
| 4.4.3. | * | Total credits (TotalCredit) | Must contain the control sum of field 4.4.4.14.5. - Credit amount (CreditAmount), excluding documents where field 4.4.4.9.1. Current document status (PaymentStatus) is of type "A". | Monetary |
| 4.4.4. | | Issued receipt document (Payment) | | N/A |
| 4.4.4.1. | * | Unique receipt identification (PaymentRefNo) | This identification is composed sequentially of the following elements: the internal code of the receipt type assigned by the application, a space, the receipt series identifier, a forward slash (/) and the sequential number of that receipt within that series. No records with the same identification can exist. The same internal document type code cannot be used in different types of documents. | Text 60 |
| 4.4.4.2. | * | Unique Document Code (ATCUD) | This field must contain the Unique Document Code. The field must be filled in with "0" (zero) until its regulation. | Text 100 |
| 4.4.4.3. | | Accounting period (Period) | The month of the taxation period from "1" to "12" must be indicated, counted from the date of its beginning. | Integer |
| 4.4.4.4. | ** | Transaction identifier (TransactionID) | Filling is mandatory in the case of an integrated accounting and billing system, even if the file type (TaxAccountingBasis) should not contain the accounting tables. The unique key from table 3. Accounting movements (GeneralLedgerEntries) where this document was posted must be indicated, respecting the rule defined there for field 3.4.3.1 - Unique accounting movement key (TransactionID). | Text 70 |
| 4.4.4.5. | * | Receipt date (TransactionDate) | Receipt issue date. | Date |
| 4.4.4.6. | * | Receipt type (PaymentType) | Must be filled in with: "RC" - Receipt issued under the Cash VAT scheme (including those relating to advances under that scheme); "RG" - Other issued receipts. | Text 2 |
| 4.4.4.7. | | Payment description (Description) | | Text 200 |
| 4.4.4.8. | | System generated number (SystemID) | Unique receipt number generated internally by the application. | Text 60 |
| 4.4.4.9. | * | Document status (DocumentStatus) | | N/A |
| 4.4.4.9.1. | * | Current receipt status (PaymentStatus) | Must be filled in with: "N" - Normal receipt and in force; "A" - Canceled receipt. | Text 1 |
| 4.4.4.9.2. | * | Receipt current status date and time (PaymentStatusDate) | Date of the last status recording to the second. Date and time type: "YYYY-MM-DDThh:mm:ss". | Date and Time |
| 4.4.4.9.3. | | Receipt status change reason (Reason) | The reason that led to the change in receipt status must be indicated. | Text 50 |
| 4.4.4.9.4. | * | User code (SourceID) | User responsible for the current status of the receipt. | Text 30 |
| 4.4.4.9.5. | * | Document source (SourcePayment) | Must be filled in with: "P" - Receipt produced in the application; "I" - Receipt integrated and produced in another application; "M" - Receipt from recovery or manual issuance. | Text 1 |
| 4.4.4.10. | | Payment method (PaymentMethod) | The payment method used must be indicated. In case of mixed payments, the amounts by type and payment date must be indicated. If there is a need to make more than one reference, this structure can be generated as many times as necessary. | N/A |
| 4.4.4.10.1. | | Payment mechanism (PaymentMechanism) | Must be filled in with: "CC" - Credit card; "CD" - Debit card; "CH" - Bank check; "CI" - International documentary credit; "CO" - Gift check or card; "CS" - Current account balance compensation; "DE" - Electronic money, for example resident in loyalty or points cards; "LC" - Commercial bill; "MB" - Payment references for Multibanco; "NU" - Cash; "OU" - Other means not indicated here; "PR" - Goods exchange; "TB" - Bank transfer or authorized direct debit; "TR" - Extra-salary compensation titles regardless of their support, for example, meal vouchers, education, etc. | Text 2 |
| 4.4.4.10.2. | * | Payment amount (PaymentAmount) | The amount by payment method must be indicated. | Monetary |
| 4.4.4.10.3. | * | Payment date (PaymentDate) | | Date |
| 4.4.4.11. | * | User code (SourceID) | User who generated the document. | Text 30 |
| 4.4.4.12. | * | Receipt recording date (SystemEntryDate) | Date of record recording to the second. Date and time type: "YYYY-MM-DDThh:mm:ss". | Date and Time |
| 4.4.4.13. | * | Customer identifier (CustomerID) | Unique key from table 2.2. - Customer table (Customer) respecting the rule defined there for field 2.2.1. - Unique customer identifier (CustomerID). | Text 30 |
| 4.4.4.14. | * | Line (Line) | | N/A |
| 4.4.4.14.1. | * | Line number (LineNumber) | Lines must be exported in the same order as they appear in the original receipt. | Integer |
| 4.4.4.14.2. | * | Source document reference (SourceDocumentID) | If there is a need to make more than one reference, this structure can be generated as many times as necessary. In case the application is integrated, the numbering structure of the source field must be used. | N/A |
| 4.4.4.14.2.1. | * | Source document number (OriginatingON) | The type, series and number of the invoice or invoice rectifying document to which the payment refers must be indicated. If the referenced document is contained in the SAF-T(PT), the numbering structure of field 4.1.4.1. - Unique sales document identification (InvoiceNo) of Table 4.1. Commercial documents to customers (SalesInvoices) must be used. | Text 60 |
| 4.4.4.14.2.2. | * | Invoice date (InvoiceDate) | The date of the invoice or invoice rectifying document to which the payment refers must be indicated. | Date |
| 4.4.4.14.2.3. | | Line description (Description) | Description of the receipt line. | Text 200 |
| 4.4.4.14.3. | | Line settlement amount (SettlementAmount) | Discounts granted when paying this document. | Monetary |
| 4.4.4.14.4. | ** | Debit amount (DebitAmount) | Receipt value of the rectifying document, without taxes and eventual discounts. | Monetary |
| 4.4.4.14.5. | ** | Credit amount (CreditAmount) | Receipt value of the invoice or rectifying document, without taxes and eventual discounts. | Monetary |
| 4.4.4.14.6. | ** | Tax rate (Tax) | In Cash VAT scheme receipts, a line must be indicated for each different VAT rate that appears in the respective line. This structure must also be generated in any other type of receipts that discriminate the types of taxes described in field 4.4.4.14.6.1. - Tax type code (TaxType). | N/A |
| 4.4.4.14.6.1. | * | Tax type code (TaxType) | In this field, the type of tax must be indicated. Must be filled in with: "VAT" - Value Added Tax, for the Cash VAT scheme; "IS" - Stamp Duty; "NS" - Not subject to VAT or IS. | Text 3 |
| 4.4.4.14.6.2. | * | Tax country/region (TaxCountryRegion) | Must be filled in according to ISO 3166-1-alpha-2 standard. In the case of the Autonomous Regions of Madeira and Azores, must be filled in with: "PT-AC" - Fiscal space of the Autonomous Region of the Azores; "PT-MA" - Fiscal space of the Autonomous Region of Madeira. | Text 5 |
| 4.4.4.14.6.3. | * | Tax code (TaxCode) | Tax code in the tax table. In case field 4.4.4.14.6.1. - Tax type code (TaxType) = VAT, must be filled in with: "RED" - Reduced rate; "INT" - Intermediate rate; "NOR" - Normal rate; "ISE" - Exempt; "OUT" - Others, applicable for special VAT regimes. In case field 4.4.4.14.6.1. - Tax type code (TaxType) = IS, must be filled in with: - The respective item code; - "ISE" - Exempt. In case of non-applicability of tax must be filled in with "NA". | Text 10 |
| 4.4.4.14.6.4. | ** | Tax rate percentage (TaxPercentage) | Filling is mandatory in case it is a tax percentage. The rate percentage corresponds to the tax applicable to field 4.4.4.14.4. - Debit amount (DebitAmount) or field 4.4.4.14.5. - Credit amount (CreditAmount). In case of exemption or non-subjection to tax, use this field with the value "0" (zero). | Decimal |
| 4.4.4.14.6.5. | ** | Tax amount (TaxAmount) | Filling is mandatory in case it is a fixed stamp duty unit amount. This value, multiplied by the quantity (Quantity) contributes to the tax payable amount (TaxPayable). | Monetary |
| 4.4.4.14.7. | ** | Tax exemption reason (TaxExemptionReason) | Filling is mandatory when fields 4.4.4.14.6.4. - Tax rate percentage (TaxPercentage) or 4.4.4.14.6.5. - Tax amount (TaxAmount) are equal to zero. The applicable legal provision must be referenced. This field must also be filled in in cases of non-subjection to the taxes referred to in table 2.5. - Tax table (TaxTable). | Text 60 |
| 4.4.4.14.8. | ** | Tax exemption code (TaxExemptionCode) | Must be filled in with the exemption or non-assessment reason code, which appears in the Software Integration Manual - Communication of Invoices to AT. | Text 3 |
| 4.4.4.15. | * | Document totals (DocumentTotals) | | N/A |
| 4.4.4.15.1. | * | Tax payable (TaxPayable) | | Monetary |
| 4.4.4.15.2. | * | Total document without taxes (NetTotal) | This field must not include the tax portions contained in table 2.5. - Tax table (TaxTable). | Monetary |
| 4.4.4.15.3. | * | Total document with taxes (GrossTotal) | This field must not reflect any withholding taxes contained in structure 4.4.4.16. - Withholding tax (WithholdingTax). | Monetary |
| 4.4.4.15.4. | | Settlement (Settlement) | Payment agreements or methods. | N/A |
| 4.4.4.15.4.1. | * | Settlement amount (SettlementAmount) | Total discounts granted when making this payment. | Monetary |
| 4.4.4.15.5. | | Currency (Currency) | Must not be generated in case the document issuance is in euro. | N/A |
| 4.4.4.15.5.1. | * | Currency code (CurrencyCode) | In case of foreign currency must be filled in according to ISO 4217 standard. | Text 3 |
| 4.4.4.15.5.2. | * | Currency amount (CurrencyAmount) | Value of field 4.4.4.15.3. - Total document with taxes (GrossTotal) in the original currency of the document. | Monetary |
| 4.4.4.15.5.3. | * | Exchange rate (ExchangeRate) | The exchange rate used in the conversion to EUR must be indicated. | Decimal |
| 4.4.4.16. | | Withholding tax (WithholdingTax) | If there is a need to make more than one reference, this structure can be generated as many times as necessary. | N/A |
| 4.4.4.16.1. | | Withholding tax type (WithholdingTaxType) | In this field, the type of tax withheld must be indicated, filling it with: "IRS" - Personal income tax; "IRC" - Corporate income tax; "IS" - Stamp duty. | Text 3 |
| 4.4.4.16.2. | | Withholding tax description (WithholdingTaxDescription) | The applicable legal provision must be indicated. In case of Tax type code (WithholdingTaxType) = IS, must be filled in with the respective item code. | Text 60 |
| 4.4.4.16.3. | * | Withholding tax amount (WithholdingTaxAmount) | The amount of tax withheld must be indicated. | Monetary |

## Annex II

### Taxonomy S - SNC base and International Accounting Standards

(referred to in Article 4 of this ordinance)

| Taxonomy Code | SNC base Code | Complete Description | Notes |
|:------------:|:-------------:|:--------------------|:------|
| 1 | 11 | Cash | |
| 2 | 12 | Demand deposits | |
| 3 | 13 | Other bank deposits | |
| 4 | 1411 | Other financial instruments - Derivatives - Potentially favorable | |
| 5 | 1412 | Other financial instruments - Derivatives - Potentially unfavorable | |
| 6 | 1421 | Other financial instruments - Financial instruments held for trading - Financial assets | |
| 7 | 1422 | Other financial instruments - Financial instruments held for trading - Financial liabilities | |
| 8 | 1431 | Other financial instruments - Other financial assets and liabilities - Other financial assets | |
| 9 | 1432 | Other financial instruments - Other financial assets and liabilities - Other financial liabilities | |
| 10 | 2111 | Customers - Customers c/c - General customers | Include accounts 2117|2118|2119 |
| 11 | 2112 | Customers - Customers c/c - Parent company customers | |
| 12 | 2113 | Customers - Customers c/c - Subsidiary company customers | |
| 13 | 2114 | Customers - Customers c/c - Associated company customers | |
| 14 | 2115 | Customers - Customers c/c - Joint venture customers | |
| 15 | 2116 | Customers - Customers c/c - Other related party customers | Include accounts 2127|2128|2129 |
| 16 | 2121 | Customers - Customers - bills receivable - General customers | |
| 17 | 2122 | Customers - Customers - bills receivable - Parent company customers | |
| 18 | 2123 | Customers - Customers - bills receivable - Subsidiary company customers | |
| 19 | 2124 | Customers - Customers - bills receivable - Associated company customers | |
| 20 | 2125 | Customers - Customers - bills receivable - Joint venture customers | |
| 21 | 2126 | Customers - Customers - bills receivable - Other related party customers | |
| 22 | 213/214/215/216/217/218/219 | Customers - Other Customers | |
| 23 | 218 | Customers - Customer advances | |
| 24 | 219 | Customers - Accumulated impairment losses - Customers c/c - General customers | |
| 25 | 219 | Customers - Accumulated impairment losses - Customers c/c - Parent company customers | |
| 26 | 219 | Customers - Accumulated impairment losses - Customers c/c - Subsidiary company customers | |
| 27 | 219 | Customers - Accumulated impairment losses - Customers c/c - Associated company customers | |
| 28 | 219 | Customers - Accumulated impairment losses - Customers c/c - Joint venture customers | |
| 29 | 219 | Customers - Accumulated impairment losses - Customers c/c - Other related party customers | |
| 30 | 219 | Customers - Accumulated impairment losses - Customers - bills receivable - General customers | |
| 31 | 219 | Customers - Accumulated impairment losses - Customers - bills receivable - Parent company customers | |
| 32 | 219 | Customers - Accumulated impairment losses - Customers - bills receivable - Subsidiary company customers | |
| 33 | 219 | Customers - Accumulated impairment losses - Customers - bills receivable - Associated company customers | |
| 34 | 219 | Customers - Accumulated impairment losses - Customers - bills receivable - Joint venture customers | |
| 35 | 219 | Customers - Accumulated impairment losses - Customers - bills receivable - Other related party customers | |
| 36 | 219 | Customers - Accumulated impairment losses - Other Customers | |
| 37 | 2211 | Suppliers - Suppliers c/c - General suppliers | Include accounts 2217|2218|2219 |
| 38 | 2212 | Suppliers - Suppliers c/c - Parent company suppliers | |
| 39 | 2213 | Suppliers - Suppliers c/c - Subsidiary company suppliers | |
| 40 | 2214 | Suppliers - Suppliers c/c - Associated company suppliers | |
| 41 | 2215 | Suppliers - Suppliers c/c - Joint venture suppliers | |
| 42 | 2216 | Suppliers - Suppliers c/c - Other related party suppliers | Include accounts 2227|2228|2229 |
| 43 | 2221 | Suppliers - Suppliers - bills payable - General suppliers | |
| 44 | 2222 | Suppliers - Suppliers - bills payable - Parent company suppliers | |
| 45 | 2223 | Suppliers - Suppliers - bills payable - Subsidiary company suppliers | |
| 46 | 2224 | Suppliers - Suppliers - bills payable - Associated company suppliers | |
| 47 | 2225 | Suppliers - Suppliers - bills payable - Joint venture suppliers | |
| 48 | 2226 | Suppliers - Suppliers - bills payable - Other related party suppliers | |
| 49 | 223/224/225/226/227/228/229 | Suppliers - Other Suppliers | |
| 50 | 228 | Suppliers - Supplier advances | |
| 51 | 229 | Suppliers - Accumulated impairment losses | |
| 52 | 23 | Personnel | |
| 53 | 24 | State and other public entities | |
| 54 | 25 | Financing obtained | |
| 55 | 26 | Shareholders/partners | |
| 56 | 27 | Other accounts receivable and payable | |
| 57 | 28 | Deferrals | |
| 58 | 29 | Provisions | |
| 59 | 31 | Purchases | |
| 60 | 32 | Inventories - Goods | |
| 61 | 33 | Raw materials, subsidiaries and consumables | |
| 62 | 34 | Finished and intermediate products | |
| 63 | 35 | By-products, waste and scrap | |
| 64 | 36 | Products and work in progress | |
| 65 | 37 | Biological assets | |
| 66 | 38 | Reclassification and regularization of inventories | |
| 67 | 39 | Inventory impairment losses | |
| 68 | 41 | Investment in financial assets | |
| 69 | 42 | Investment properties | |
| 70 | 43 | Tangible fixed assets | |
| 71 | 44 | Intangible assets | |
| 72 | 45 | Investments in progress | |
| 73 | 46 | Non-current assets held for sale | |
| 74 | 51 | Capital | |
| 75 | 52 | Own shares | |
| 76 | 53 | Other equity instruments | |
| 77 | 54 | Share premiums | |
| 78 | 55 | Reserves | |
| 79 | 56 | Retained earnings | |
| 80 | 57 | Adjustments to financial assets | |
| 81 | 58 | Revaluation surplus | |
| 82 | 59 | Other changes in equity | |
| 83 | 61 | Cost of goods sold and materials consumed | |
| 84 | 62 | External supplies and services | |
| 85 | 63 | Personnel expenses | |
| 86 | 64 | Depreciation and amortization expenses | |
| 87 | 65 | Impairment losses | |
| 88 | 66 | Losses on reduction of fair value | |
| 89 | 67 | Provisions for the period | |
| 90 | 68 | Other expenses and losses | |
| 91 | 69 | Financing expenses and losses | |
| 92 | 71 | Sales | |
| 93 | 72 | Services rendered | |
| 94 | 73 | Changes in inventories of production | |
| 95 | 74 | Work for the entity itself | |
| 96 | 75 | Operating subsidies | |
| 97 | 76 | Reversals | |
| 98 | 77 | Gains on increase of fair value | |
| 99 | 78 | Other income and gains | |
| 100 | 79 | Interest, dividends and other similar income | |
| 101 | 81 | Income tax | |
| 102 | 818 | Net income for the period | |
| 103 | 89 | Dividends advanced | |

## Annex III

### Taxonomy M - SNC Microentities

(referred to in Article 4 of this ordinance)

| Taxonomy Code | SNC Microentities Code | Complete Description | Notes |
|:------------:|:---------------------:|:--------------------|:------|
| 1 | 11 | Cash | |
| 2 | 12 | Demand deposits | |
| 3 | 13 | Other bank deposits | |
| 4 | 14 | Other financial instruments | |
| 5 | 21 | Customers | |
| 6 | 22 | Suppliers | |
| 7 | 23 | Personnel | |
| 8 | 24 | State and other public entities | |
| 9 | 25 | Financing obtained | |
| 10 | 26 | Shareholders/partners | |
| 11 | 27 | Other accounts receivable and payable | |
| 12 | 28 | Deferrals | |
| 13 | 29 | Provisions | |
| 14 | 31 | Purchases | |
| 15 | 32 | Inventories - Goods | |
| 16 | 33 | Raw materials, subsidiaries and consumables | |
| 17 | 34 | Finished and intermediate products | |
| 18 | 35 | By-products, waste and scrap | |
| 19 | 36 | Products and work in progress | |
| 20 | 37 | Biological assets | |
| 21 | 38 | Reclassification and regularization of inventories | |
| 22 | 39 | Inventory impairment losses | |
| 23 | 41 | Investment in financial assets | |
| 24 | 43 | Tangible fixed assets | |
| 25 | 44 | Intangible assets | |
| 26 | 45 | Investments in progress | |
| 27 | 46 | Non-current assets held for sale | |
| 28 | 51 | Capital | |
| 29 | 52 | Own shares | |
| 30 | 53 | Other equity instruments | |
| 31 | 54 | Share premiums | |
| 32 | 55 | Reserves | |
| 33 | 56 | Retained earnings | |
| 34 | 57 | Adjustments to financial assets | |
| 35 | 58 | Revaluation surplus | |
| 36 | 59 | Other changes in equity | |
| 37 | 61 | Cost of goods sold and materials consumed | |
| 38 | 62 | External supplies and services | |
| 39 | 63 | Personnel expenses | |
| 40 | 64 | Depreciation and amortization expenses | |
| 41 | 65 | Impairment losses | |
| 42 | 67 | Provisions for the period | |
| 43 | 68 | Other expenses and losses | |
| 44 | 69 | Financing expenses and losses | |
| 45 | 71 | Sales | |
| 46 | 72 | Services rendered | |
| 47 | 73 | Changes in inventories of production | |
| 48 | 74 | Work for the entity itself | |
| 49 | 75 | Operating subsidies | |
| 50 | 76 | Reversals | |
| 51 | 78 | Other income and gains | |
| 52 | 79 | Interest, dividends and other similar income | |
| 53 | 81 | Income tax | |
| 54 | 818 | Net income for the period | |
| 55 | 89 | Dividends advanced | |
