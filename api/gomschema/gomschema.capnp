@0xdc8319c2c83c62ef;  # unique ID

enum AllegianceType {
    allegNone @0;
    allegAlliance @1;
    allegEmpire @2;
    allegFederation @3;
    allegIndependent @4;
    allegPilotsFederation @5;
}

enum FacilityType {
    fTNone @0;
    fTCivilianOutpost @1;
    fTCommercialOutpost @2;
    fTCoriolisStarport @3;
    fTIndustrialOutpost @4;
    fTMilitaryOutpost @5;
    fTMiningOutpost @6;
    fTOcellusStarport @7;
    fTOrbisStarport @8;
    fTScientificOutpost @9;
    fTPlanetaryOutpost @10;
    fTPlanetaryPort @11;
    fTPlanetarySettlement @12;
    fTMegaship @13;
    fTAsteroidBase @14;
    fTFleetCarrier @15;
}

enum FeatureBits {
	fBitZero @0;
    fBitMarket @1;
    fBitBlackMarket @2;
    fBitCommodities @3;
    fBitDocking @4;
    fBitFleet @5;
    fBitLargePad @6;
    fBitMediumPad @7;
    fBitOutfitting @8;
    fBitPlanetary @9;
    fBitRearm @10;
    fBitRefuel @11;
    fBitRepair @12;
    fBitShipyard @13;
    fBitSmallPad @14;
}

enum GovernmentType {
    govNone @0;
    govAnarchy @1;
    govCommunism @2;
    govConfederacy @3;
    govCooperative @4;
    govCorporate @5;
    govDemocracy @6;
    govDictatorship @7;
    govFeudal @8;
    govPatronage @9;
    govPrison @10;
    govPrisonColony @11;
    govTheocracy @12;
}

enum PadSize {
    padNone @0;
    padSmall @1;
    padMedium @2;
    padLarge @3;
}

enum SecurityLevel {
    securityNone @0;
    securityAnarchy @1;
    securityLow @2;
    securityMedium @3;
    securityHigh @4;
}

struct Coordinate {
    x @0 :Float64;
    y @1 :Float64;
    z @2 :Float64;
}

struct Datum {
    union {
        commodity @0 :Commodity;
        system @1 :System;
        facility @2 :Facility;
        facilityListing @3 :FacilityListing;
    }
}

struct UserData {
    identity @0 :Text;      # What this data is for.
    data @1 :Data;          # Opaque data.
}

struct Commodity {
    enum Category {
        catNone @0;
        catChemicals @1;
        catConsumerItems @2;
        catLegalDrugs @3;
        catFoods @4;
        catIndustrialMaterials @5;
        catMachinery @6;
        catMedicines @7;
        catMetals @8;
        catMinerals @9;
        catSlavery @10;
        catTechnology @11;
        catTextiles @12;
        catWaste @13;
        catWeapons @14;
        catUnknown @15;
        catSalvage @16;
    }

    id @0 :UInt32;
    name @1 :Text;
    timestampUTC @2 :UInt64 = 0;
    categoryId @3 :Category;
    isRare @4 :Bool = false;
    isNonMarketable @5 :Bool = false;
    averageCr @6 :UInt32;
}

struct System {
    id @0 :UInt32;
    name @1 :Text;
    timestampUTC @2 :UInt64 = 0;
    position @3 :Coordinate;
    populated @4 :Bool = true;
    needsPermit @5 :Bool = false;
    securityLevel @6 :SecurityLevel = securityMedium;
    government @7 :GovernmentType = govCorporate;
    allegiance @8 :AllegianceType = allegIndependent;
}

struct Facility {
    id @0 :UInt32;
    systemId @1 :UInt32;
    name @2 :Text;
    timestampUTC @3 :UInt64 = 0;
    facilityType @4 :FacilityType;
    features @5 :UInt16;
    lsFromStar @6 :UInt32;
    government @7 :GovernmentType = govCorporate;
    allegiance @8 :AllegianceType = allegIndependent;
}

struct CommodityListing {
    commodityId @0 :UInt32;
    supplyUnits @1 :UInt32;
    supplyCredits @2 :UInt32;
    demandUnits @3 :UInt32;
    demandCredits @4 :UInt32;
    timestampUTC @5 :UInt64;
}

struct FacilityListing {
    id @0 :UInt32;
    listings @1 :List(CommodityListing);
}

struct Header {
    enum Type {
        tInvalid    @0;
        tHeader     @1;
        tCommodity @2;
        tSystem    @3;
        tFacility  @4;
        tListing   @5;
    }

    headerType @0 :Type;
    data @1 :List(Datum);
    userData @2 :List(UserData) = [];
}
