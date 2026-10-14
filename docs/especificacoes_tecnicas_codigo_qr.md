# ESPECIFICAÇÕES TÉCNICAS Código de Barras Bidimensional Código QR

Portaria n. ${ }^{\circ}$ 195 / 2020
Versão 1.1

# HISTÓRICO DE ALTERAÇÕES

| DATA | ALTERAÇÕES |
| :-- | :-- |
| 13-08-2020 | Criação do documento |
| 30-10-2020 | - Atualização das regras de composição da mensagem a constar no <br> código QR e das instruções dos campos; <br> - Alteração da abrangência do código "L" <br> - Novos exemplos |
|  |  |

# ÍNDICE

1. Introdução ..... 3
2. Código de barras bidimensional (código QR) ..... 4
3. Composição da mensagem a constar no código QR ..... 4
4. Estrutura e conteúdos do código QR ..... 5
5. Exemplos de construção do código QR ..... 14
5.1. Exemplo 1 - Fatura ..... 14
5.2. Exemplo 2 - Fatura simplificada ..... 15
5.3. Exemplo 3 - Fatura pró-forma ..... 16
5.4. Exemplo 4 - Documento de transporte ..... 17
5.5. Exemplo 5 - Fatura com taxa de imposto estrangeiro ..... 18
5.6. Exemplo 6 - Nota de débito retificativa de imposto ..... 19
5.7. Exemplo 7 - Fatura com IVA - Regime da Margem ..... 20

# 1. Introdução

O Decreto-Lei n. ${ }^{\circ}$ 28/2019, de 15 de fevereiro, veio consolidar e modernizar as normas relativas à faturação, introduzindo aspetos inovadores como o código único de documento e o código de barras bidimensional (código QR), que visam a simplificação na comunicação de faturas por parte de pessoas singulares para determinação das respetivas despesas dedutíveis em sede de IRS, incrementando, simultaneamente, o controlo das operações realizadas pelos sujeitos passivos, tendo em vista combater a economia informal, a fraude e a evasão fiscal.

O código de barras bidimensional (código QR) deve constar, obrigatoriamente, nas faturas e outros documentos fiscalmente relevantes, emitidos por programas certificados pela Autoridade Tributária e Aduaneira (AT), nos termos do artigo 4. ${ }^{\circ}$ do Decreto-Lei n. ${ }^{\circ}$ 28/2019, de 15 de fevereiro.

O presente documento define as especificações técnicas para a geração do código de barras bidimensional (código QR), previstas Portaria n. ${ }^{\circ}$ 195/2020, de 13 de agosto.

Os produtores de programas de faturação certificados pela AT, devem garantir a correta geração do código QR, respeitando as especificações técnicas constantes no presente documento, devendo ainda garantir a sua perfeita legibilidade dentro do corpo do documento, independentemente do suporte em que seja apresentado ao cliente.

Em documentos com mais do que uma página, o código de barras bidimensional (código QR) pode constar na primeira ou na última página.

Para apoiar os produtores de programas de faturação na implementação do código QR nos seus programas, são identificados no presente documento os campos que devem integrar na mensagem, com as respetivas instruções de preenchimento, bem como um exemplo para cada campo.

São ainda disponibilizados exemplos de construção de código QR, contemplando diversos tipos de documentos e diferentes mensagens.

# 2. Código de barras bidimensional (código QR)

A elaboração do código de barras bidimensional (código QR) deve obedecer às seguintes especificações:
(a) Taxa de Recuperação de Erro (ECC): "M";
(b) Tipo: Byte;
(c) Pontos por módulo (Size): 2;
(d) Versão: $v=9$ (valor mínimo);
(e) Dimensões de imagem: mínimo $30 	imes 30$ milímetros;
(f) Margem de Segurança (Margin): $0,25 \mathrm{~cm}$.

## 3. Composição da mensagem a constar no código QR

Na criação da mensagem a incorporar no código QR devem ser observadas as seguintes regras:
(a) Cada campo será formado pela concatenação do valor da coluna "Código", constante na tabela do ponto 4, «:» (dois pontos) e o respetivo valor da coluna "Descrição", sem espaços;
(b) Os campos assim criados, e sempre pela ordem indicada na tabela do ponto 4, deverão ser concatenados com o separador «*» (asterisco);
(c) Os campos monetários deverão ser representados em euros, com «.» (ponto) como separador decimal e sempre com duas (2) casas decimais, independentemente do número de casas decimais (inferior ou superior) apresentado na base de dados e na exportação para o ficheiro SAF-T (PT). Nos documentos emitidos em moeda diferente de Euro, os montantes apresentados deverão ser previamente convertidos em euros;

(d) Os campos assinalados com «+» são de criação obrigatória;
(e) Os campos assinalados com «++», são opcionais, mas deverão ser criados sempre que para eles exista informação;
(f) Nos campos opcionais, na ausência de informação não deverá ser criado o respetivo campo;
(g) Na criação dos campos com os códigos I1 a I8, J1 a J8 e K1 a K8, representativos dos espaços fiscais para efeitos de IVA (por exemplo, PT, PTAC e PT-MA), terá sempre de existir, pelo menos, um espaço fiscal, até ao máximo de três espaços fiscais em simultâneo, nacionais ou estrangeiros;
(h) Na composição do campo com o código I1, no caso de documento emitido sem indicação da taxa de IVA, que deva constar na tabela 4.2, 4.3 ou 4.4 do SAF-T (PT), deverá ser preenchido com «I1»«:»«0»;
(i) Nenhum valor da coluna "Descrição" poderá ultrapassar o tamanho máximo definido na tabela do ponto 4 ;
(j) Na composição do campo com o código "S", sempre que necessário, os elementos que o compõem serão concatenados com «;» (ponto e vírgula) sem espaços;
(k) Para as diferentes taxas de IVA existentes no documento (isento, reduzida, intermédia e normal) devem constar, nos respetivos campos, os totais acumulados de base tributável e IVA.

# 4. Estrutura e conteúdos do código QR

Na geração do código QR, deverá ser respeitada a estrutura definida na tabela seguinte.

| Código | Descrição | Tamanho Máximo | Instruções | Exemplo | Obrigatório |
| :--: | :--: | :--: | :--: | :--: | :--: |
| A | NIF do <br> emitente | 9 | Preencher com o NIF do <br> emitente sem espaços e sem <br> prefixo do País, de acordo <br> com o campo <br> TaxRegistrationNumber do <br> SAF-T (PT). | A:123456789 | $+$ |
|  | Separador | 1 |  | * | $+$ |
| B | NIF do <br> adquirente | 30 | Preencher com o NIF do <br> adquirente sem prefixo do <br> país, de acordo com o campo <br> CustomerTaxID do SAF-T <br> (PT). Na emissão de <br> documento a "Consumidor <br> Final" preencher com <br> 999999990. | B:999999990 | $+$ |
|  | Separador | 1 |  | * | $+$ |
| C | País do <br> adquirente | 12 | Preencher de acordo com o <br> campo Country da tabela de <br> clientes do SAF-T (PT). | C:PT | $+$ |
|  | Separador | 1 |  | * | $+$ |
| D | Tipo de <br> documento | 2 | Preencher de acordo com a <br> tipologia do SAF-T (PT) - <br> campos InvoiceType, <br> MovementType, WorkType <br> ou PaymentType. | D:FT | $+$ |
|  | Separador | 1 |  | * | $+$ |
| E | Estado do <br> documento | 1 | Preencher de acordo com a <br> tipologia do SAF-T (PT) - <br> campos InvoiceStatus, <br> MovementStatus, WorkStatus <br> ou PaymentStatus. | E:N | $+$ |
|  | Separador | 1 |  | * | $+$ |
| F | Data do <br> documento | 8 | Utilizar o formato <br> YYYYMMDD. <br> Corresponde aos campos <br> InvoiceDate, MovementDate, <br> WorkDate ou <br> TransactionDate do SAF-T <br> (PT) sem hífens. | F:20191231 | $+$ |
|  | Separador | 1 |  | * | $+$ |

| Código | Descrição | Tamanho <br> Máximo | Instruções | Exemplo | Obrigatório |
| :--: | :--: | :--: | :--: | :--: | :--: |
| G | Identificação única do documento | 60 | Preencher de acordo com a tipologia do SAF-T (PT) campos InvoiceNo, DocumentNumber ou PaymentRefNo. | G:FT <br> AB2019/0035 | $+$ |
|  | Separador | 1 |  | * | $+$ |
| H | ATCUD | 70 | Preencher com o código único de documento, de acordo com os campos ATCUD do SAF-T (PT). | H:CSDF7T5H0035 | $+$ |
|  | Separador | 1 |  | * | $+$ |
| I1 | Espaço fiscal | 5 | Preencher de acordo com as notas técnicas do campo TaxCountryRegion do SAF-T (PT). <br> No caso de documento sem indicação da taxa de IVA, que deva constar na tabela 4.2, 4.3 ou 4.4 do SAF-T (PT), preencher com «0» (I1:0). | II:PT | $+$ |
|  | Separador | 1 |  | * | $++$ |
| I2 | Base tributável isenta de IVA | 16 | Valor total da base tributável isenta de IVA, incluindo as operações tributáveis em sede do Imposto do Selo (sejam isentas ou não de Imposto do Selo). <br> Formatar com duas casas decimais, com "." como separador decimal e sem separador de milhares. | I2:12000.00 | $++$ |
|  | Separador | 1 |  | * | $++$ |
| I3 | Base tributável de IVA à taxa reduzida | 16 | Valor total da base tributável sujeita à taxa reduzida de IVA. <br> Formatar com duas casas decimais, com "." como separador decimal e sem separador de milhares. | I3:15000.00 | $++$ |
|  | Separador | 1 |  | * | $++$ |

| Código | Descrição | Tamanho Máximo | Instruções | Exemplo | Obrigatório |
| :--: | :--: | :--: | :--: | :--: | :--: |
| 14 | Total de IVA à taxa reduzida | 16 | Valor total do IVA à taxa reduzida no documento. <br> Formatar com duas casas decimais, com "." como separador decimal e sem separador de milhares. | I4:900.00 | $++$ |
|  | Separador | 1 |  | * | $++$ |
| 15 | Base tributável de IVA à taxa intermédia | 16 | Valor total da base tributável sujeita à taxa intermédia de IVA. Formatar com duas casas decimais, com "." como separador decimal e sem separador de milhares. | I5:50000.00 | $++$ |
|  | Separador | 1 |  | * | $++$ |
| 16 | Total de IVA à taxa intermédia | 16 | Valor total do IVA à taxa intermédia no documento. <br> Formatar com duas casas decimais, com "." como separador decimal e sem separador de milhares. | I6:6500.00 | $++$ |
|  | Separador | 1 |  | * | $++$ |
| 17 | Base tributável de IVA à taxa normal | 16 | Valor total da base tributável sujeita à taxa normal de IVA. <br> Formatar com duas casas decimais, com "." como separador decimal e sem separador de milhares. | I7:80000.00 | $++$ |
|  | Separador | 1 |  | * | $++$ |
| 18 | Total de IVA à taxa normal | 16 | Valor total do IVA à taxa normal no documento. <br> Formatar com duas casas decimais, com "." como separador decimal e sem separador de milhares. | I8:18400.00 | $++$ |
|  | Separador | 1 |  | * | $++$ |
| J1 | Espaço fiscal | 5 | Preencher de acordo com as notas técnicas do campo TaxCountryRegion do SAF-T (PT). | J1:PT-AC | $++$ |
|  | Separador | 1 |  | * | $++$ |

| Código | Descrição | Tamanho <br> Máximo | Instruções | Exemplo | Obrigatório |
| :--: | :--: | :--: | :--: | :--: | :--: |
| J2 | Base tributável isenta de IVA | 16 | Valor total da base tributável isenta de IVA, incluindo as operações tributáveis em sede do Imposto do Selo (sejam isentas ou não de Imposto do Selo). <br> Formatar com duas casas decimais, com "." como separador decimal e sem separador de milhares. | J2:10000.00 | $++$ |
|  | Separador | 1 |  | * | $++$ |
| J3 | Base tributável de IVA à taxa reduzida | 16 | Valor total da base tributável sujeita à taxa reduzida de IVA. <br> Formatar com duas casas decimais, com "." como separador decimal e sem separador de milhares. | J3:25000.56 | $++$ |
|  | Separador | 1 |  | * | $++$ |
| J4 | Total de IVA à taxa reduzida | 16 | Valor total do IVA à taxa reduzida no documento. <br> Formatar com duas casas decimais, com "." como separador decimal e sem separador de milhares. | J4:1000.02 | $++$ |
|  | Separador | 1 |  | * | $++$ |
| J5 | Base tributável de IVA à taxa intermédia | 16 | Valor total da base tributável sujeita à taxa intermédia de IVA. Formatar com duas casas decimais, com "." como separador decimal e sem separador de milhares. | J5:75000.00 | $++$ |
|  | Separador | 1 |  | * | $++$ |
| J6 | Total de IVA à taxa intermédia | 16 | Valor total do IVA à taxa intermédia no documento. <br> Formatar com duas casas decimais, com "." como separador decimal e sem separador de milhares. | J6:6750.00 | $++$ |
|  | Separador | 1 |  | * | $++$ |

| Código | Descrição | Tamanho <br> Máximo | Instruções | Exemplo | Obrigatório |
| :--: | :--: | :--: | :--: | :--: | :--: |
| J7 | Base tributável de IVA à taxa normal | 16 | Valor total da base tributável sujeita à taxa normal de IVA. Formatar com duas casas decimais, com "." como separador decimal e sem separador de milhares. | J7:100000.00 | $++$ |
|  | Separador | 1 |  | * | $++$ |
| J8 | Total de IVA à taxa normal | 16 | Valor total do IVA à taxa normal no documento. <br> Formatar com duas casas decimais, com "." como separador decimal e sem separador de milhares. | J8:18000.00 | $++$ |
|  | Separador | 1 |  | * | $++$ |
| K1 | Espaço fiscal | 5 | Preencher de acordo com as notas técnicas do campo TaxCountryRegion do SAF-T (PT). | K1:PT-MA | $++$ |
|  | Separador | 1 |  | * | $++$ |
| K2 | Base tributável isenta de IVA | 16 | Valor total da base tributável isenta de IVA, incluindo as operações tributáveis em sede do Imposto do Selo (sejam isentas ou não de Imposto do Selo). <br> Formatar com duas casas decimais, com "." como separador decimal e sem separador de milhares. | K2:5000.00 | $++$ |
|  | Separador | 1 |  | * | $++$ |
| K3 | Base tributável de IVA à taxa reduzida | 16 | Valor total da base tributável sujeita à taxa reduzida de IVA. <br> Formatar com duas casas decimais, com "." como separador decimal e sem separador de milhares. | K3:12500.00 | $++$ |
|  | Separador | 1 |  | * | $++$ |

| Código | Descrição | Tamanho <br> Máximo | Instruções | Exemplo | Obrigatório |
| :--: | :--: | :--: | :--: | :--: | :--: |
| K4 | Total de IVA à taxa reduzida | 16 | Valor total do IVA à taxa reduzida no documento. Formatar com duas casas decimais, com "." como separador decimal e sem separador de milhares. | K4:625.00 | $++$ |
|  | Separador | 1 |  | * | $++$ |
| K5 | Base tributável de IVA à taxa intermédia | 16 | Valor total da base tributável sujeita à taxa intermédia de IVA. <br> Formatar com duas casas decimais, com "." como separador decimal e sem separador de milhares. | K5:25000.00 | $++$ |
|  | Separador | 1 |  | * | $++$ |
| K6 | Total de IVA à taxa intermédia | 16 | Valor total do IVA à taxa intermédia no documento. <br> Formatar com duas casas decimais, com "." como separador decimal e sem separador de milhares. | K6:3000.00 | $++$ |
|  | Separador | 1 |  | * | $++$ |
| K7 | Base tributável de IVA à taxa normal | 16 | Valor total da base tributável sujeita à taxa normal de IVA. <br> Formatar com duas casas decimais, com "." como separador decimal e sem separador de milhares. | K7:40000.00 | $++$ |
|  | Separador | 1 |  | * | $++$ |
| K8 | Total de IVA à taxa normal | 16 | Valor total do IVA à taxa normal no documento. <br> Formatar com duas casas decimais, com "." como separador decimal e sem separador de milhares. | K8:8800.00 | $++$ |
|  | Separador | 1 |  | * | $++$ |

| Código | Descrição | Tamanho <br> Máximo | Instruções | Exemplo | Obrigatório |
| :--: | :--: | :--: | :--: | :--: | :--: |
| L | Não sujeito / não tributável em IVA / outras situações | 16 | Valor total das operações não sujeitas / não tributáveis em IVA e das operações sujeitas a IVA cujo documento não discrimine ou não possa discriminar o valor do imposto (ex.: documentos de transporte e recibos gerais sem discriminação de imposto ou regimes especiais de tributação da margem Agências de viagens, bens em segunda mão, objetos de arte, objetos de coleção e antiguidades). <br> Formatar com duas casas decimais, com "." como separador decimal e sem separador de milhares. | L:100.00 | $++$ |
|  | Separador | 1 |  | * | $++$ |
| M | Imposto do <br> Selo | 16 | Valor total do Imposto do Selo no documento. <br> Formatar com duas casas decimais, com "." como separador decimal e sem separador de milhares. | M:25.00 | $++$ |
|  | Separador | 1 |  | * | $+$ |
| N | Total de impostos | 16 | Valor total de IVA e Imposto do Selo - campo TaxPayable do SAF-T (PT). <br> Formatar com duas casas decimais, com "." como separador decimal e sem separador de milhares. | N:64000.02 | $+$ |
|  | Separador | 1 |  | * | $+$ |
| 0 | Total do documento com impostos | 16 | Valor do total do documento campo GrossTotal do SAF-T (PT). <br> Formatar com duas casas decimais, com "." como separador decimal e sem separador de milhares. | O:513600.58 | $+$ |
|  | Separador | 1 |  | * | $++$ |

| Código | Descrição | Tamanho <br> Máximo | Instruções | Exemplo | Obrigatório |
| :--: | :--: | :--: | :--: | :--: | :--: |
| P | Retenções na fonte | 16 | Valor do total das retenções na fonte - campo WithholdingTaxAmount do SAF-T (PT). <br> Formatar com duas casas decimais, com "." como separador decimal e sem separador de milhares. | P:100.00 | $++$ |
|  | Separador | 1 |  | * | $+$ |
| Q | 4 carateres do Hash | 4 | Preencher de acordo com alínea a) do n. ${ }^{\circ}$ 3 do art. ${ }^{\circ}$ 6 .${ }^{\circ}$ da Portaria n. ${ }^{\circ}$ 363/2010, de 23 de junho. | Q:kLp0 | $+$ |
|  | Separador | 1 |  | * | $+$ |
| R | № do certificado | 4 | Preencher com o n. ${ }^{\circ}$ do certificado atribuído pela AT de acordo com o campo SoftwareCertificateNumber do SAF-T (PT). | R:9999 | $+$ |
|  | Separador | 1 |  | * | $++$ |
| S | Outras informações | 65 | Campo de preenchimento livre, em que pode, por exemplo, ser indicada informação para pagamento (ex.: de IBAN ou Ref. ${ }^{9}$ MB, com o separador «;»). <br> Este campo não pode conter o caracter asterisco (*). | S:TB;PT00000 000000000000 000000;51350 0.58 <br> (ou <br> S:MB;entidade ;referência;val or) | $++$ |

Nota: As remissões a "SAF-T (PT)" devem ser entendidas como referências à estrutura de dados a que se refere o n. ${ }^{\circ}$ 1 .${ }^{\circ}$ da Portaria n. ${ }^{\circ}$ 321-A / 2007, de 26 de março.

# 5. Exemplos de construção do código QR

### 5.1. Exemplo 1 - Fatura

- Mensagem com os elementos da tabela acima, com indicação do IBAN e valores de IVA nos espaços fiscais PT, PT-AC e PT-MA:

A:123456789*B:999999990*C:PT*D:FT*E:N*F:20191231*G:FT AB2019/0035*H:CSDF7T5H-0035*I1:PT*I2:12000.00*I3:15000.00*I4:900.00*I5:50000.00*I6:6500.00*I7:80 000.00*I8:18400.00*J1:PT-AC*J2:10000.00*J3:25000.56*J4:1000.02*J5:75000.00*J6:6750.00*J7:10000 $0.00 *$ J8:18000.00*K1:PT-MA*K2:5000.00*K3:12500.00*K4:625.00*K5:25000.00*K6:3000.00*K7:40000. 00*K8:8800.00*L:100.00*M:25.00*N:64000.02*O:513600.58*P:100.00*Q:kLp 0*R:9999*S:TB;PT000000000000000000000000;513500.58

![QR Code Example 1](https://via.placeholder.com/150)

Figura 1 - Exemplo de código QR [47,096mm] e [452 caracteres]
(Size $=2, \mathrm{ECC}=\mathrm{M}$, Margin $=2(0,25 \mathrm{~cm})$, Module[versão]=17->85x85)

### 5.2. Exemplo 2 - Fatura simplificada

- Fatura simplificada emitida com os campos da tabela abaixo:

| Código | Descrição | Exemplo | Obrigatório |
| :--: | :-- | :--: | :--: |
| A | NIF do emitente | A:123456789 | + |
| B | NIF do adquirente | B:999999990 | + |
| C | País do adquirente | C:PT | + |
| D | Tipo de documento | D:FS | + |
| E | Estado do documento | E:N | + |
| F | Data do documento | F:20190812 | + |
| G | Identificação única do documento | G:FS CDVF/12345 | + |
| H | ATCUD | H:CDF7T5HD-12345 | + |
| I1 | Espaço fiscal | I1:PT | + |
| I7 | Base tributável de IVA à taxa normal | I7:0.65 | ++ |
| I8 | Total de IVA à taxa normal | I8:0.15 | ++ |
| N | Total de impostos | N:0.15 | + |
| O | Total do documento com impostos | O:0.80 | + |
| Q | 4 carateres do Hash | Q:YhGV | + |
| R | Nº do certificado | R:9999 | + |
| S | Outras informações | S:NU;0.80 | ++ |

- Mensagem criada:

A:123456789*B:999999990*C:PT*D:FS*E:N*F:20190812*G:FS CDVF/12345*H:CDF7T5HD-12345*I1:PT*I7:0.65*I8:0.15*N:0.15*O:0.80*Q:YhGV*R:9999*S:NU;0.80

![QR Code Example 2](https://via.placeholder.com/150)

Figura 2 - Exemplo de código QR [30,162mm] e [141 caracteres]
(Size=2, ECC=M, Margin $=2(0,25 \mathrm{~cm})$, Module[versão]=9 -> 53x53)

### 5.3. Exemplo 3 - Fatura pró-forma

- Fatura pró-forma emitida com os campos da tabela abaixo:

| Código | Descrição | Exemplo | Obrigatório |
| :--: | :--: | :--: | :--: |
| A | NIF do emitente | A:500000000 | $+$ |
| B | NIF do adquirente | B:123456789 | $+$ |
| C | País do adquirente | C:PT | $+$ |
| D | Tipo de documento | D:PF | $+$ |
| E | Estado do documento | E:N | $+$ |
| F | Data do documento | F:20190123 | $+$ |
| G | Identificação única do documento | G:PF G2019CB/145789 | $+$ |
| H | ATCUD | H:HB6FT7RV-145789 | $+$ |
| 11 | Espaço fiscal | I1:PT | $+$ |
| 12 | Base tributável isenta | I2:12345.34 | $++$ |
| 13 | Base tributável de IVA à taxa reduzida | I3:12532.65 | $++$ |
| 14 | Total de IVA à taxa reduzida | I4:751.96 | $++$ |
| 15 | Base tributável de IVA à taxa intermédia | I5:52789.00 | $++$ |
| 16 | Total de IVA à taxa intermédia | I6:6862.57 | $++$ |
| 17 | Base tributável de IVA à taxa normal | I7:32425.69 | $++$ |
| 18 | Total de IVA à taxa normal | I8:7457.91 | $++$ |
| N | Total de impostos | N:15072.44 | $+$ |
| 0 | Total do documento com impostos | O:125165.12 | $+$ |
| Q | 4 carateres do Hash | Q:r/fY | $+$ |
| R | № do certificado | R:9999 | $+$ |

- Mensagem criada:

A:500000000*B:123456789*C:PT*D:PF*E:N*F:20190123*G:PF G2019CB/145789*H:HB6FT7RV-145789*I1:PT*I2:12345.34*I3:12532.65*I4:751.96*I5:52789.00*I6:6862.57*I7:32425.69*18:7457.91*N:15072.44*O:125165.12*Q:r/fY*R:9999

![QR Code Example 3](https://via.placeholder.com/150)

Figura 3 - Exemplo de código QR [32,279mm] e [209 caracteres]
(Size=2, ECC=M, Margin= $2(0,25 \mathrm{~cm})$, Module[versão]=10 -> 57x57)

### 5.4. Exemplo 4 - Documento de transporte

- Guia de transporte não valorizada, emitida com os campos da tabela abaixo:

| Código | Descrição | Exemplo | Obrigatório |
| :--: | :-- | :--: | :--: |
| A | NIF do emitente | A:500000000 | + |
| B | NIF do adquirente | B:123456789 | + |
| C | País do adquirente | C:PT | + |
| D | Tipo de documento | D:GT | + |
| E | Estado do documento | E:N | + |
| F | Data do documento | F:20190720 | + |
| G | Identificação única do documento | G:GT G234CB/50987 | + |
| H | ATCUD | H:GTVX4Y8B-50987 | + |
| I1 | Espaço fiscal | I1:0 | + |
| N | Total de impostos (*) | N:0.00 | + |
| O | Total do documento com impostos (*) | O:0.00 | + |
| Q | 4 carateres do Hash | Q:5ulg | + |
| R | Nº do certificado | R:9999 | + |

(*) Ainda que este registo se encontre valorizado na base de dados do programa, se o documento não exibir os valores, estes também não poderão constar do código QR. Igual procedimento deverá ser adotado nos documentos que sejam suscetíveis de apresentação ao cliente e que possibilitem a conferência de mercadorias ou de prestação de serviços que se encontrem valorizados na base de dados, quando os valores não constem dos documentos.

- Mensagem criada:

A:500000000*B:123456789*C:PT*D:GT*E:N*F:20190720*G:GT G234CB/50987*H:GTVX4Y8B-50987*I1:0*N:0.00*O:0.00*Q:5ulg*R:9999

![QR Code Example 4](https://via.placeholder.com/150)

Figura 4 - Exemplo de código QR [30,162mm] e [116 caracteres]
(Size=2, ECC=M, Margin= 2 (0,25cm), Module[versão]=9 -> 53x53)

### 5.5. Exemplo 5 - Fatura com taxa de imposto estrangeiro

- Fatura emitida com taxa de IVA de França, com um valor de 100 euros, sobre a qual incide IVA à taxa normal (20%):

| Código | Descrição | Exemplo | Obrigatório |
| :--: | :-- | :-- | :--: |
| A | NIF do emitente | A:123456789 | + |
| B | NIF do adquirente | B:4443332215 | + |
| C | País do adquirente | C:FR | + |
| D | Tipo de documento | D:FT | + |
| E | Estado do documento | E:N | + |
| F | Data do documento | F:20190526 | + |
| G | Identificação única do documento | G:ABC BNH/4561 | + |
| H | ATCUD | H:DK5ZJ2HN-4561 | + |
| I1 | Espaço fiscal | I1:FR | + |
| I7 | Base tributável de IVA à taxa normal | I7:100.00 | ++ |
| I8 | Total de IVA à taxa normal | I8:20.00 | ++ |
| N | Total de impostos | N:20.00 | + |
| O | Total do documento com impostos | O:120.00 | + |
| Q | 4 carateres do Hash | Q:YJRE | + |
| R | Nº do certificado | R:9999 | + |

- Mensagem criada:

A:123456789*B:4443332215*C:FR*D:FT*E:N*F:20190526*G:ABC BNH/4561*H:DK5ZJ2HN-4561*I1:FR*I7:100.00*I8:20.00*N:20.00*O:120.00*Q:YJRE*R:9999

![QR Code Example 5](https://via.placeholder.com/150)

Figura 5 - Exemplo de código QR [30,162mm] e [136 caracteres]
(Size=2, ECC=M, Margin $=2(0,25 \mathrm{~cm})$, Module[versão]=9 -> 53x53)

### 5.6. Exemplo 6 - Nota de débito retificativa de imposto

- Nota de débito emitida, apenas para retificação do imposto insuficientemente liquidado numa fatura emitida, no montante de $€ 26,50$ :

| Código | Descrição | Exemplo | Obrigatório |
| :--: | :-- | :--: | :--: |
| A | NIF do emitente | A:123456789 | + |
| B | NIF do adquirente | B:500000000 | + |
| C | País do adquirente | C:PT | + |
| D | Tipo de documento | D:ND | + |
| E | Estado do documento | E:N | + |
| F | Data do documento | F:20190216 | + |
| G | Identificação única do documento | G:M1F KLG/6145 | + |
| H | ATCUD | H:RQD8L6DG-6145 | + |
| I1 | Espaço fiscal | I1:PT-MA | + |
| I6 | Total de IVA à taxa intermédia | I6:26.50 | ++ |
| N | Total de impostos | N:26.50 | + |
| O | Total do documento com impostos | O:26.50 | + |
| Q | 4 carateres do Hash | Q:h1rB | + |
| R | Nº do certificado | R:9999 | + |

- Mensagem criada:

A:123456789*B:500000000*C:PT*D:ND*E:N*F:20190216*G: M1F KLG/6145*H: RQD8L6DG-6145*I1:PT-MA*I6:26.50*N:26.50*O:26.50*Q:h1rB*R:9999

![QR Code Example 6](https://via.placeholder.com/150)

Figura 6 - Exemplo de código QR [30,162mm] e [129 caracteres]
(Size=2, ECC=M, Margin= 2 (0,25cm), Module[versão]=9 -> 53x53)

### 5.7. Exemplo 7 - Fatura com IVA - Regime da Margem

- Fatura emitida num regime especial de tributação da margem ( $€ 1.000,00$ com imposto incluído), contendo, em simultâneo, uma transação sujeita ao regime geral de IVA ( $€ 59,00$ com imposto incluído).

| Código | Descrição | Exemplo | Obrigatório |
| :--: | :-- | :--: | :--: |
| A | NIF do emitente | A:500000000 | + |
| B | NIF do adquirente | B:123456789 | + |
| C | País do adquirente | C:PT | + |
| D | Tipo de documento | D:FT | + |
| E | Estado do documento | E:N | + |
| F | Data do documento | F:20191124 | + |
| G | Identificação única do documento | G:NF 19A/789145 | + |
| H | ATCUD | H:JL9DS4TT-789145 | + |
| I1 | Espaço fiscal | I1:PT-AC | + |
| I7 | Base tributável de IVA à taxa normal | I7:50.00 | ++ |
| I8 | Total de IVA à taxa normal | I8:9.00 | ++ |
| L | Não sujeito / não tributável em IVA / <br> outras situações | L:1000.00 | ++ |
| N | Total de impostos (*) | N:9.00 | + |
| O | Total do documento com impostos | O:1059.00 | + |
| Q | 4 carateres do Hash | Q:d8/K | + |
| R | № do certificado | R:9999 | + |

(*) Se o programa permitir determinar o montante do IVA subjacente à transação relativa ao regime da margem e este constar do seu repositório de dados, terá que ser exportado para o ficheiro SAF-T (PT). No entanto, como legalmente, os documentos emitidos não podem discriminar o IVA, este também não poderá constar do código QR.

- Mensagem criada:

A:500000000*B:123456789*C:PT*D:FT*E:N*F:20191124*G:NF 19A/789145*H:JL9DS4TT-789145*I1:PT-AC*I7:50.00*I8:9.00*L:1000.00*N:9.00*O:1059.00*Q:d8/K*R:9999

![QR Code Example 7](https://via.placeholder.com/150)

Figura 7 - Exemplo de código QR [30,162mm] e [150 caracteres]
(Size=2, ECC=M, Margin $=2(0,25 \mathrm{~cm})$, Module[versão]=9 -> 53x53)
