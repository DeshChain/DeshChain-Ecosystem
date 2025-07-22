package keeper

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitializeStandardTemplates creates all standard trade finance document templates
func (k Keeper) InitializeStandardTemplates(ctx context.Context) error {
	templates := []DocumentTemplate{
		k.createLetterOfCreditTemplate(),
		k.createBillOfLadingTemplate(),
		k.createCommercialInvoiceTemplate(),
		k.createPackingListTemplate(),
		k.createCertificateOfOriginTemplate(),
		k.createInsuranceDocumentTemplate(),
		k.createCustomsDeclarationTemplate(),
		k.createStandbyLCTemplate(),
		k.createTransferableLCTemplate(),
		k.createRevolvingLCTemplate(),
	}
	
	for _, template := range templates {
		if err := k.CreateDocumentTemplate(ctx, &template); err != nil {
			return fmt.Errorf("failed to create template %s: %w", template.TemplateName, err)
		}
	}
	
	return nil
}

// Letter of Credit Template
func (k Keeper) createLetterOfCreditTemplate() DocumentTemplate {
	return DocumentTemplate{
		TemplateID:   "LC-001",
		TemplateName: "Standard Letter of Credit",
		DocumentType: LetterOfCredit,
		Version:      "1.0",
		Fields: []TemplateField{
			{
				FieldID:   "lc_number",
				FieldName: "LC Number",
				FieldType: TextField,
				Required:  true,
				Format:    "^LC-[0-9]{10}$",
				Position:  FieldPosition{Row: 1, Column: 1},
			},
			{
				FieldID:   "issue_date",
				FieldName: "Issue Date",
				FieldType: DateField,
				Required:  true,
				Format:    "YYYY-MM-DD",
				Position:  FieldPosition{Row: 1, Column: 2},
			},
			{
				FieldID:   "expiry_date",
				FieldName: "Expiry Date",
				FieldType: DateField,
				Required:  true,
				Format:    "YYYY-MM-DD",
				Position:  FieldPosition{Row: 1, Column: 3},
				ValidationRules: []FieldValidation{
					{
						Type:    "date_after",
						Params:  map[string]string{"field": "issue_date"},
						Message: "Expiry date must be after issue date",
					},
				},
			},
			{
				FieldID:   "applicant",
				FieldName: "Applicant (Importer)",
				FieldType: TextField,
				Required:  true,
				Position:  FieldPosition{Row: 2, Column: 1, Width: 3},
				ValidationRules: []FieldValidation{
					{
						Type:    "min_length",
						Params:  map[string]string{"value": "3"},
						Message: "Applicant name must be at least 3 characters",
					},
				},
			},
			{
				FieldID:   "beneficiary",
				FieldName: "Beneficiary (Exporter)",
				FieldType: TextField,
				Required:  true,
				Position:  FieldPosition{Row: 3, Column: 1, Width: 3},
			},
			{
				FieldID:   "issuing_bank",
				FieldName: "Issuing Bank",
				FieldType: TextField,
				Required:  true,
				Position:  FieldPosition{Row: 4, Column: 1, Width: 2},
			},
			{
				FieldID:   "advising_bank",
				FieldName: "Advising Bank",
				FieldType: TextField,
				Required:  false,
				Position:  FieldPosition{Row: 4, Column: 3},
			},
			{
				FieldID:   "amount",
				FieldName: "LC Amount",
				FieldType: AmountField,
				Required:  true,
				Position:  FieldPosition{Row: 5, Column: 1},
				ValidationRules: []FieldValidation{
					{
						Type:    "min_value",
						Params:  map[string]string{"value": "0"},
						Message: "Amount must be positive",
					},
				},
			},
			{
				FieldID:   "currency",
				FieldName: "Currency",
				FieldType: TextField,
				Required:  true,
				Format:    "^[A-Z]{3}$",
				Position:  FieldPosition{Row: 5, Column: 2},
			},
			{
				FieldID:   "tolerance",
				FieldName: "Amount Tolerance (%)",
				FieldType: NumberField,
				Required:  false,
				DefaultValue: 5,
				Position:  FieldPosition{Row: 5, Column: 3},
			},
			{
				FieldID:   "goods_description",
				FieldName: "Description of Goods",
				FieldType: TextField,
				Required:  true,
				Position:  FieldPosition{Row: 6, Column: 1, Width: 3, Height: 3},
			},
			{
				FieldID:   "incoterms",
				FieldName: "Incoterms",
				FieldType: TextField,
				Required:  true,
				Format:    "^(EXW|FCA|CPT|CIP|DAP|DPU|DDP|FAS|FOB|CFR|CIF)$",
				Position:  FieldPosition{Row: 9, Column: 1},
			},
			{
				FieldID:   "port_of_loading",
				FieldName: "Port of Loading",
				FieldType: TextField,
				Required:  true,
				Position:  FieldPosition{Row: 9, Column: 2},
			},
			{
				FieldID:   "port_of_discharge",
				FieldName: "Port of Discharge",
				FieldType: TextField,
				Required:  true,
				Position:  FieldPosition{Row: 9, Column: 3},
			},
			{
				FieldID:   "latest_shipment_date",
				FieldName: "Latest Shipment Date",
				FieldType: DateField,
				Required:  true,
				Position:  FieldPosition{Row: 10, Column: 1},
			},
			{
				FieldID:   "presentation_period",
				FieldName: "Presentation Period (days)",
				FieldType: NumberField,
				Required:  true,
				DefaultValue: 21,
				Position:  FieldPosition{Row: 10, Column: 2},
			},
			{
				FieldID:   "documents_required",
				FieldName: "Documents Required",
				FieldType: TableField,
				Required:  true,
				Position:  FieldPosition{Row: 11, Column: 1, Width: 3, Height: 5},
			},
			{
				FieldID:   "additional_conditions",
				FieldName: "Additional Conditions",
				FieldType: TextField,
				Required:  false,
				Position:  FieldPosition{Row: 16, Column: 1, Width: 3, Height: 3},
			},
			{
				FieldID:   "charges",
				FieldName: "Charges",
				FieldType: TextField,
				Required:  true,
				Position:  FieldPosition{Row: 19, Column: 1},
			},
			{
				FieldID:   "confirmation_required",
				FieldName: "Confirmation Required",
				FieldType: TextField,
				Required:  true,
				DefaultValue: "WITHOUT",
				Position:  FieldPosition{Row: 19, Column: 2},
			},
			{
				FieldID:   "transferable",
				FieldName: "Transferable",
				FieldType: TextField,
				Required:  true,
				DefaultValue: "NO",
				Position:  FieldPosition{Row: 19, Column: 3},
			},
		},
		ValidationRules: []ValidationRule{
			{
				RuleID:      "lc_validity",
				RuleName:    "LC Validity Check",
				Description: "Ensure LC validity period is reasonable",
				Condition:   "expiry_date - issue_date <= 180 days",
				ErrorMessage: "LC validity period cannot exceed 180 days",
			},
			{
				RuleID:      "shipment_before_expiry",
				RuleName:    "Shipment Before Expiry",
				Description: "Latest shipment date must be before expiry",
				Condition:   "latest_shipment_date < expiry_date",
				ErrorMessage: "Latest shipment date must be before LC expiry",
			},
		},
		Layout: LayoutSpecification{
			PageSize:    "A4",
			Orientation: "Portrait",
			Margins:     Margins{Top: 20, Bottom: 20, Left: 25, Right: 25},
			FontFamily:  "Arial",
			FontSize:    11,
		},
		Metadata: TemplateMetadata{
			Category:        "Trade Finance",
			Jurisdiction:    "International",
			ComplianceRules: []string{"UCP 600", "ISBP 745"},
			SupportedFormats: []string{"PDF", "JSON", "XML"},
		},
	}
}

// Bill of Lading Template
func (k Keeper) createBillOfLadingTemplate() DocumentTemplate {
	return DocumentTemplate{
		TemplateID:   "BL-001",
		TemplateName: "Ocean Bill of Lading",
		DocumentType: BillOfLading,
		Version:      "1.0",
		Fields: []TemplateField{
			{
				FieldID:   "bl_number",
				FieldName: "B/L Number",
				FieldType: TextField,
				Required:  true,
				Format:    "^[A-Z0-9-]+$",
				Position:  FieldPosition{Row: 1, Column: 1},
			},
			{
				FieldID:   "shipper",
				FieldName: "Shipper",
				FieldType: AddressField,
				Required:  true,
				Position:  FieldPosition{Row: 2, Column: 1, Width: 2, Height: 3},
			},
			{
				FieldID:   "consignee",
				FieldName: "Consignee",
				FieldType: AddressField,
				Required:  true,
				Position:  FieldPosition{Row: 2, Column: 3, Height: 3},
			},
			{
				FieldID:   "notify_party",
				FieldName: "Notify Party",
				FieldType: AddressField,
				Required:  false,
				Position:  FieldPosition{Row: 5, Column: 1, Width: 3, Height: 2},
			},
			{
				FieldID:   "vessel_name",
				FieldName: "Vessel Name",
				FieldType: TextField,
				Required:  true,
				Position:  FieldPosition{Row: 7, Column: 1},
			},
			{
				FieldID:   "voyage_number",
				FieldName: "Voyage Number",
				FieldType: TextField,
				Required:  true,
				Position:  FieldPosition{Row: 7, Column: 2},
			},
			{
				FieldID:   "port_of_loading",
				FieldName: "Port of Loading",
				FieldType: TextField,
				Required:  true,
				Position:  FieldPosition{Row: 8, Column: 1},
			},
			{
				FieldID:   "port_of_discharge",
				FieldName: "Port of Discharge",
				FieldType: TextField,
				Required:  true,
				Position:  FieldPosition{Row: 8, Column: 2},
			},
			{
				FieldID:   "place_of_delivery",
				FieldName: "Place of Delivery",
				FieldType: TextField,
				Required:  false,
				Position:  FieldPosition{Row: 8, Column: 3},
			},
			{
				FieldID:   "marks_numbers",
				FieldName: "Marks & Numbers",
				FieldType: TextField,
				Required:  false,
				Position:  FieldPosition{Row: 9, Column: 1, Height: 3},
			},
			{
				FieldID:   "cargo_description",
				FieldName: "Description of Cargo",
				FieldType: TextField,
				Required:  true,
				Position:  FieldPosition{Row: 9, Column: 2, Height: 3},
			},
			{
				FieldID:   "packages_quantity",
				FieldName: "Number of Packages",
				FieldType: NumberField,
				Required:  true,
				Position:  FieldPosition{Row: 9, Column: 3},
			},
			{
				FieldID:   "gross_weight",
				FieldName: "Gross Weight",
				FieldType: NumberField,
				Required:  true,
				Position:  FieldPosition{Row: 10, Column: 3},
			},
			{
				FieldID:   "measurement",
				FieldName: "Measurement",
				FieldType: NumberField,
				Required:  false,
				Position:  FieldPosition{Row: 11, Column: 3},
			},
			{
				FieldID:   "freight_terms",
				FieldName: "Freight Terms",
				FieldType: TextField,
				Required:  true,
				Format:    "^(PREPAID|COLLECT)$",
				Position:  FieldPosition{Row: 12, Column: 1},
			},
			{
				FieldID:   "freight_amount",
				FieldName: "Freight Amount",
				FieldType: AmountField,
				Required:  false,
				Position:  FieldPosition{Row: 12, Column: 2},
			},
			{
				FieldID:   "place_of_issue",
				FieldName: "Place of Issue",
				FieldType: TextField,
				Required:  true,
				Position:  FieldPosition{Row: 13, Column: 1},
			},
			{
				FieldID:   "date_of_issue",
				FieldName: "Date of Issue",
				FieldType: DateField,
				Required:  true,
				Position:  FieldPosition{Row: 13, Column: 2},
			},
			{
				FieldID:   "number_of_originals",
				FieldName: "Number of Original B/Ls",
				FieldType: NumberField,
				Required:  true,
				DefaultValue: 3,
				Position:  FieldPosition{Row: 13, Column: 3},
			},
			{
				FieldID:   "carrier_signature",
				FieldName: "Carrier Signature",
				FieldType: SignatureField,
				Required:  true,
				Position:  FieldPosition{Row: 14, Column: 1, Width: 3},
			},
		},
		ValidationRules: []ValidationRule{
			{
				RuleID:      "bl_date_validity",
				RuleName:    "B/L Date Validity",
				Description: "B/L issue date cannot be future dated",
				Condition:   "date_of_issue <= today",
				ErrorMessage: "Bill of Lading cannot be future dated",
			},
			{
				RuleID:      "weight_required",
				RuleName:    "Weight Requirement",
				Description: "Gross weight must be specified",
				Condition:   "gross_weight > 0",
				ErrorMessage: "Gross weight must be greater than zero",
			},
		},
		Layout: LayoutSpecification{
			PageSize:    "A4",
			Orientation: "Portrait",
			Margins:     Margins{Top: 15, Bottom: 15, Left: 20, Right: 20},
			FontFamily:  "Courier New",
			FontSize:    10,
		},
		Metadata: TemplateMetadata{
			Category:        "Shipping Documents",
			Jurisdiction:    "International Maritime",
			ComplianceRules: []string{"Hague-Visby Rules", "Hamburg Rules"},
			SupportedFormats: []string{"PDF", "EDI", "JSON"},
		},
	}
}

// Commercial Invoice Template
func (k Keeper) createCommercialInvoiceTemplate() DocumentTemplate {
	return DocumentTemplate{
		TemplateID:   "INV-001",
		TemplateName: "Commercial Invoice",
		DocumentType: Invoice,
		Version:      "1.0",
		Fields: []TemplateField{
			{
				FieldID:   "invoice_number",
				FieldName: "Invoice Number",
				FieldType: TextField,
				Required:  true,
				Position:  FieldPosition{Row: 1, Column: 1},
			},
			{
				FieldID:   "invoice_date",
				FieldName: "Invoice Date",
				FieldType: DateField,
				Required:  true,
				Position:  FieldPosition{Row: 1, Column: 2},
			},
			{
				FieldID:   "seller_details",
				FieldName: "Seller/Exporter",
				FieldType: AddressField,
				Required:  true,
				Position:  FieldPosition{Row: 2, Column: 1, Width: 2, Height: 3},
			},
			{
				FieldID:   "buyer_details",
				FieldName: "Buyer/Importer",
				FieldType: AddressField,
				Required:  true,
				Position:  FieldPosition{Row: 2, Column: 3, Height: 3},
			},
			{
				FieldID:   "reference_numbers",
				FieldName: "Reference Numbers",
				FieldType: TableField,
				Required:  false,
				Position:  FieldPosition{Row: 5, Column: 1, Width: 3},
			},
			{
				FieldID:   "items_table",
				FieldName: "Line Items",
				FieldType: TableField,
				Required:  true,
				Position:  FieldPosition{Row: 6, Column: 1, Width: 3, Height: 10},
				ValidationRules: []FieldValidation{
					{
						Type:    "table_columns",
						Params:  map[string]string{"columns": "item_code,description,quantity,unit_price,total_price"},
						Message: "Items table must have required columns",
					},
				},
			},
			{
				FieldID:   "subtotal",
				FieldName: "Subtotal",
				FieldType: AmountField,
				Required:  true,
				Position:  FieldPosition{Row: 16, Column: 3},
			},
			{
				FieldID:   "tax_amount",
				FieldName: "Tax Amount",
				FieldType: AmountField,
				Required:  false,
				DefaultValue: 0,
				Position:  FieldPosition{Row: 17, Column: 3},
			},
			{
				FieldID:   "shipping_charges",
				FieldName: "Shipping Charges",
				FieldType: AmountField,
				Required:  false,
				DefaultValue: 0,
				Position:  FieldPosition{Row: 18, Column: 3},
			},
			{
				FieldID:   "total_amount",
				FieldName: "Total Amount",
				FieldType: AmountField,
				Required:  true,
				Position:  FieldPosition{Row: 19, Column: 3},
				ValidationRules: []FieldValidation{
					{
						Type:    "calculated_field",
						Params:  map[string]string{"formula": "subtotal + tax_amount + shipping_charges"},
						Message: "Total amount calculation error",
					},
				},
			},
			{
				FieldID:   "currency",
				FieldName: "Currency",
				FieldType: TextField,
				Required:  true,
				Format:    "^[A-Z]{3}$",
				Position:  FieldPosition{Row: 19, Column: 2},
			},
			{
				FieldID:   "payment_terms",
				FieldName: "Payment Terms",
				FieldType: TextField,
				Required:  true,
				Position:  FieldPosition{Row: 20, Column: 1, Width: 3},
			},
			{
				FieldID:   "incoterms",
				FieldName: "Incoterms",
				FieldType: TextField,
				Required:  true,
				Position:  FieldPosition{Row: 21, Column: 1},
			},
			{
				FieldID:   "country_of_origin",
				FieldName: "Country of Origin",
				FieldType: TextField,
				Required:  true,
				Position:  FieldPosition{Row: 21, Column: 2},
			},
			{
				FieldID:   "authorized_signature",
				FieldName: "Authorized Signature",
				FieldType: SignatureField,
				Required:  true,
				Position:  FieldPosition{Row: 22, Column: 1, Width: 2},
			},
		},
		ValidationRules: []ValidationRule{
			{
				RuleID:      "amount_positive",
				RuleName:    "Positive Amounts",
				Description: "All amounts must be positive",
				Condition:   "subtotal >= 0 && total_amount > 0",
				ErrorMessage: "Invoice amounts must be positive",
			},
		},
		Layout: LayoutSpecification{
			PageSize:    "A4",
			Orientation: "Portrait",
			Margins:     Margins{Top: 20, Bottom: 20, Left: 25, Right: 25},
			FontFamily:  "Arial",
			FontSize:    11,
		},
		Metadata: TemplateMetadata{
			Category:        "Commercial Documents",
			Jurisdiction:    "International",
			ComplianceRules: []string{"ICC Rules", "WTO Standards"},
			SupportedFormats: []string{"PDF", "XML", "JSON", "CSV"},
		},
	}
}

// Helper types for templates
type FieldPosition struct {
	Row    int
	Column int
	Width  int // Grid columns to span
	Height int // Grid rows to span
}

type FieldValidation struct {
	Type    string
	Params  map[string]string
	Message string
}

type ValidationRule struct {
	RuleID       string
	RuleName     string
	Description  string
	Condition    string
	ErrorMessage string
}

type LayoutSpecification struct {
	PageSize    string
	Orientation string
	Margins     Margins
	FontFamily  string
	FontSize    int
}

type Margins struct {
	Top    int
	Bottom int
	Left   int
	Right  int
}

type TemplateMetadata struct {
	Category         string
	Jurisdiction     string
	ComplianceRules  []string
	SupportedFormats []string
	CustomFields     map[string]interface{}
}

// Advanced LC Types Templates

// Standby Letter of Credit Template
func (k Keeper) createStandbyLCTemplate() DocumentTemplate {
	baseLC := k.createLetterOfCreditTemplate()
	standbyLC := baseLC
	standbyLC.TemplateID = "SBLC-001"
	standbyLC.TemplateName = "Standby Letter of Credit"
	
	// Add standby-specific fields
	standbyFields := []TemplateField{
		{
			FieldID:   "standby_type",
			FieldName: "Standby Type",
			FieldType: TextField,
			Required:  true,
			Position:  FieldPosition{Row: 20, Column: 1},
			ValidationRules: []FieldValidation{
				{
					Type:    "enum",
					Params:  map[string]string{"values": "Performance,Financial,Direct Pay,Counter"},
					Message: "Invalid standby LC type",
				},
			},
		},
		{
			FieldID:   "default_statement",
			FieldName: "Statement of Default",
			FieldType: TextField,
			Required:  true,
			Position:  FieldPosition{Row: 21, Column: 1, Width: 3, Height: 3},
		},
		{
			FieldID:   "drawing_conditions",
			FieldName: "Drawing Conditions",
			FieldType: TextField,
			Required:  true,
			Position:  FieldPosition{Row: 24, Column: 1, Width: 3, Height: 3},
		},
	}
	
	standbyLC.Fields = append(standbyLC.Fields, standbyFields...)
	standbyLC.Metadata.ComplianceRules = append(standbyLC.Metadata.ComplianceRules, "ISP98")
	
	return standbyLC
}

// Transferable Letter of Credit Template
func (k Keeper) createTransferableLCTemplate() DocumentTemplate {
	baseLC := k.createLetterOfCreditTemplate()
	transferableLC := baseLC
	transferableLC.TemplateID = "TLC-001"
	transferableLC.TemplateName = "Transferable Letter of Credit"
	
	// Update transferable field
	for i, field := range transferableLC.Fields {
		if field.FieldID == "transferable" {
			transferableLC.Fields[i].DefaultValue = "YES"
			break
		}
	}
	
	// Add transfer-specific fields
	transferFields := []TemplateField{
		{
			FieldID:   "first_beneficiary",
			FieldName: "First Beneficiary",
			FieldType: TextField,
			Required:  true,
			Position:  FieldPosition{Row: 25, Column: 1, Width: 3},
		},
		{
			FieldID:   "second_beneficiary",
			FieldName: "Second Beneficiary",
			FieldType: TextField,
			Required:  false,
			Position:  FieldPosition{Row: 26, Column: 1, Width: 3},
		},
		{
			FieldID:   "transfer_bank",
			FieldName: "Transferring Bank",
			FieldType: TextField,
			Required:  false,
			Position:  FieldPosition{Row: 27, Column: 1, Width: 2},
		},
		{
			FieldID:   "transfer_conditions",
			FieldName: "Transfer Conditions",
			FieldType: TextField,
			Required:  true,
			Position:  FieldPosition{Row: 28, Column: 1, Width: 3, Height: 2},
		},
		{
			FieldID:   "transfer_charges",
			FieldName: "Transfer Charges Borne By",
			FieldType: TextField,
			Required:  true,
			DefaultValue: "First Beneficiary",
			Position:  FieldPosition{Row: 30, Column: 1},
		},
	}
	
	transferableLC.Fields = append(transferableLC.Fields, transferFields...)
	
	// Add transfer-specific validation
	transferableLC.ValidationRules = append(transferableLC.ValidationRules, ValidationRule{
		RuleID:      "transfer_amount_check",
		RuleName:    "Transfer Amount Validation",
		Description: "Transfer amount cannot exceed original LC amount",
		Condition:   "transfer_amount <= original_amount",
		ErrorMessage: "Transfer amount exceeds original LC amount",
	})
	
	return transferableLC
}

// Revolving Letter of Credit Template
func (k Keeper) createRevolvingLCTemplate() DocumentTemplate {
	baseLC := k.createLetterOfCreditTemplate()
	revolvingLC := baseLC
	revolvingLC.TemplateID = "RLC-001"
	revolvingLC.TemplateName = "Revolving Letter of Credit"
	
	// Add revolving-specific fields
	revolvingFields := []TemplateField{
		{
			FieldID:   "revolving_type",
			FieldName: "Revolving Type",
			FieldType: TextField,
			Required:  true,
			Position:  FieldPosition{Row: 25, Column: 1},
			ValidationRules: []FieldValidation{
				{
					Type:    "enum",
					Params:  map[string]string{"values": "Automatic,Non-Automatic"},
					Message: "Invalid revolving type",
				},
			},
		},
		{
			FieldID:   "revolving_basis",
			FieldName: "Revolving Basis",
			FieldType: TextField,
			Required:  true,
			Position:  FieldPosition{Row: 25, Column: 2},
			ValidationRules: []FieldValidation{
				{
					Type:    "enum",
					Params:  map[string]string{"values": "Time,Value,Cumulative,Non-Cumulative"},
					Message: "Invalid revolving basis",
				},
			},
		},
		{
			FieldID:   "revolving_frequency",
			FieldName: "Revolving Frequency",
			FieldType: TextField,
			Required:  true,
			Position:  FieldPosition{Row: 25, Column: 3},
		},
		{
			FieldID:   "number_of_revolutions",
			FieldName: "Number of Revolutions",
			FieldType: NumberField,
			Required:  true,
			Position:  FieldPosition{Row: 26, Column: 1},
		},
		{
			FieldID:   "revolving_amount",
			FieldName: "Amount per Revolution",
			FieldType: AmountField,
			Required:  true,
			Position:  FieldPosition{Row: 26, Column: 2},
		},
		{
			FieldID:   "total_revolving_amount",
			FieldName: "Total Revolving Amount",
			FieldType: AmountField,
			Required:  true,
			Position:  FieldPosition{Row: 26, Column: 3},
			ValidationRules: []FieldValidation{
				{
					Type:    "calculated_field",
					Params:  map[string]string{"formula": "number_of_revolutions * revolving_amount"},
					Message: "Total revolving amount calculation error",
				},
			},
		},
		{
			FieldID:   "reinstatement_conditions",
			FieldName: "Reinstatement Conditions",
			FieldType: TextField,
			Required:  true,
			Position:  FieldPosition{Row: 27, Column: 1, Width: 3, Height: 2},
		},
	}
	
	revolvingLC.Fields = append(revolvingLC.Fields, revolvingFields...)
	
	// Add revolving-specific validation
	revolvingLC.ValidationRules = append(revolvingLC.ValidationRules, ValidationRule{
		RuleID:      "revolving_period_check",
		RuleName:    "Revolving Period Validation",
		Description: "All revolutions must complete before LC expiry",
		Condition:   "last_revolution_date < expiry_date",
		ErrorMessage: "Revolving period extends beyond LC expiry",
	})
	
	return revolvingLC
}

// Additional standard templates (simplified for brevity)

func (k Keeper) createPackingListTemplate() DocumentTemplate {
	return DocumentTemplate{
		TemplateID:   "PL-001",
		TemplateName: "Packing List",
		DocumentType: PackingList,
		Version:      "1.0",
		Fields: []TemplateField{
			{
				FieldID:   "packing_list_number",
				FieldName: "Packing List Number",
				FieldType: TextField,
				Required:  true,
				Position:  FieldPosition{Row: 1, Column: 1},
			},
			{
				FieldID:   "date",
				FieldName: "Date",
				FieldType: DateField,
				Required:  true,
				Position:  FieldPosition{Row: 1, Column: 2},
			},
			{
				FieldID:   "invoice_reference",
				FieldName: "Invoice Reference",
				FieldType: TextField,
				Required:  true,
				Position:  FieldPosition{Row: 1, Column: 3},
			},
			// Additional fields omitted for brevity
		},
		Layout: LayoutSpecification{
			PageSize:    "A4",
			Orientation: "Portrait",
			Margins:     Margins{Top: 20, Bottom: 20, Left: 25, Right: 25},
			FontFamily:  "Arial",
			FontSize:    10,
		},
		Metadata: TemplateMetadata{
			Category:        "Shipping Documents",
			Jurisdiction:    "International",
			ComplianceRules: []string{"ICC Standards"},
			SupportedFormats: []string{"PDF", "Excel", "JSON"},
		},
	}
}

func (k Keeper) createCertificateOfOriginTemplate() DocumentTemplate {
	return DocumentTemplate{
		TemplateID:   "COO-001",
		TemplateName: "Certificate of Origin",
		DocumentType: CertificateOfOrigin,
		Version:      "1.0",
		Fields: []TemplateField{
			{
				FieldID:   "certificate_number",
				FieldName: "Certificate Number",
				FieldType: TextField,
				Required:  true,
				Position:  FieldPosition{Row: 1, Column: 1},
			},
			{
				FieldID:   "issue_date",
				FieldName: "Issue Date",
				FieldType: DateField,
				Required:  true,
				Position:  FieldPosition{Row: 1, Column: 2},
			},
			{
				FieldID:   "issuing_authority",
				FieldName: "Issuing Authority",
				FieldType: TextField,
				Required:  true,
				Position:  FieldPosition{Row: 2, Column: 1, Width: 3},
			},
			// Additional fields omitted for brevity
		},
		Layout: LayoutSpecification{
			PageSize:    "A4",
			Orientation: "Portrait",
			Margins:     Margins{Top: 20, Bottom: 20, Left: 25, Right: 25},
			FontFamily:  "Arial",
			FontSize:    11,
		},
		Metadata: TemplateMetadata{
			Category:        "Certificates",
			Jurisdiction:    "International",
			ComplianceRules: []string{"WTO Rules of Origin"},
			SupportedFormats: []string{"PDF", "JSON"},
		},
	}
}

func (k Keeper) createInsuranceDocumentTemplate() DocumentTemplate {
	return DocumentTemplate{
		TemplateID:   "INS-001",
		TemplateName: "Marine Insurance Certificate",
		DocumentType: InsuranceDocument,
		Version:      "1.0",
		Fields: []TemplateField{
			{
				FieldID:   "policy_number",
				FieldName: "Policy Number",
				FieldType: TextField,
				Required:  true,
				Position:  FieldPosition{Row: 1, Column: 1},
			},
			{
				FieldID:   "certificate_number",
				FieldName: "Certificate Number",
				FieldType: TextField,
				Required:  true,
				Position:  FieldPosition{Row: 1, Column: 2},
			},
			{
				FieldID:   "issue_date",
				FieldName: "Issue Date",
				FieldType: DateField,
				Required:  true,
				Position:  FieldPosition{Row: 1, Column: 3},
			},
			// Additional fields omitted for brevity
		},
		Layout: LayoutSpecification{
			PageSize:    "A4",
			Orientation: "Portrait",
			Margins:     Margins{Top: 20, Bottom: 20, Left: 25, Right: 25},
			FontFamily:  "Arial",
			FontSize:    11,
		},
		Metadata: TemplateMetadata{
			Category:        "Insurance Documents",
			Jurisdiction:    "International",
			ComplianceRules: []string{"Institute Cargo Clauses"},
			SupportedFormats: []string{"PDF", "JSON"},
		},
	}
}

func (k Keeper) createCustomsDeclarationTemplate() DocumentTemplate {
	return DocumentTemplate{
		TemplateID:   "CD-001",
		TemplateName: "Customs Declaration",
		DocumentType: CustomsDeclaration,
		Version:      "1.0",
		Fields: []TemplateField{
			{
				FieldID:   "declaration_number",
				FieldName: "Declaration Number",
				FieldType: TextField,
				Required:  true,
				Position:  FieldPosition{Row: 1, Column: 1},
			},
			{
				FieldID:   "declaration_date",
				FieldName: "Declaration Date",
				FieldType: DateField,
				Required:  true,
				Position:  FieldPosition{Row: 1, Column: 2},
			},
			{
				FieldID:   "customs_office",
				FieldName: "Customs Office",
				FieldType: TextField,
				Required:  true,
				Position:  FieldPosition{Row: 1, Column: 3},
			},
			// Additional fields omitted for brevity
		},
		Layout: LayoutSpecification{
			PageSize:    "A4",
			Orientation: "Portrait",
			Margins:     Margins{Top: 20, Bottom: 20, Left: 25, Right: 25},
			FontFamily:  "Arial",
			FontSize:    10,
		},
		Metadata: TemplateMetadata{
			Category:        "Customs Documents",
			Jurisdiction:    "Country-Specific",
			ComplianceRules: []string{"WCO Standards", "National Customs Laws"},
			SupportedFormats: []string{"PDF", "XML", "EDI"},
		},
	}
}