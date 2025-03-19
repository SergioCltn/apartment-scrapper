import logging
import re
from dataclasses import asdict, dataclass, fields, is_dataclass
from typing import Any, Self

import utils

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s - %(levelname)s - %(message)s",
)


class FlatSerializable:
    def to_flat_dict(self, prefix: str = "") -> dict[str, Any]:
        """To be implemented by child classes"""
        raise NotImplementedError(
            "Subclasses must implement to_flat_dict()",
        )


class BaseFromRaw:
    @classmethod
    def from_raw(cls, raw: str) -> Self:
        """To be implemented by child classes"""
        raise NotImplementedError(
            "Subclasses must implement from_raw()",
        )


@dataclass
class Address(FlatSerializable):
    street: str | None = None
    number: str | None = None
    neighborhood: str | None = None
    district: str | None = None
    city: str | None = None
    region: str | None = None
    autonomous_community: str | None = None

    @classmethod
    def from_location(cls, location: str) -> Self:
        """Create an Address instance from a location string.

        Args:
            location (str): The location string to parse.

        Returns:
            Address: An Address instance with parsed data.

        """
        address = cls()
        parts = [p.strip() for p in location.split(",")]

        match = re.match(r"^(.*?)(?:[\s,-]+(\d+|s/n))?$", parts[0])
        if match:
            street = match.group(1).strip()
            address.street = (
                None if street.startswith("Barrio") else street
            )
            address.number = match.group(2) if match.group(2) else None

        # Handle addresses that might have number split to second field
        if not address.number and re.match(r"^-?\d+|s/n$", parts[1]):
            address.number = parts[1]

        for part in parts:
            if part.startswith("Distrito"):
                address.district = part
            elif part.startswith("Comarca"):
                address.region = part
            elif part.startswith("Barrio"):
                address.neighborhood = part

        address.city = parts[-3]
        address.region = parts[-2]
        address.autonomous_community = parts[-1]

        return address

    def to_flat_dict(self, prefix="address") -> dict[str, Any]:
        return {f"{prefix}_{k}": v for k, v in asdict(self).items()}


@dataclass
class BasicFeatures(FlatSerializable, BaseFromRaw):
    n_rooms: int | None = None
    sqm_constructed: float | None = None
    sqm_usable: float | None = None
    bathrooms: int | None = None
    terrace: bool = False
    state: str | None = None
    orientation: str | None = None
    balcony: bool = False
    built_in: int | None = None
    heating: str | None = None
    reduced_mobility: str | None = None
    garage: str | None = None

    storage_room: bool = False
    built_in_wardrobes: bool = False
    air_conditioning: bool = False
    semi_detached_house: str | None = None
    green_areas: bool = False
    pool: bool = False
    new_building_development: bool = False

    @classmethod
    def from_raw(cls, raw: str) -> Self:
        basic_features = cls()
        for v in raw.split(";"):
            if "m² construidos" in v or "m² útiles" in v:
                match = re.search(
                    r"(?:(\d+)\s*m² construidos)?(?:,\s*)?(?:(\d+)\s*m² útiles)?",
                    v,
                )
                if match is not None:
                    constructed, usable = match.groups()
                    if constructed:
                        basic_features.sqm_constructed = float(
                            constructed,
                        )
                    if usable:
                        basic_features.sqm_usable = float(usable)

            if "habitación" in v or "habitaciones" in v:
                if "Sin " in v:
                    basic_features.n_rooms = 0
                else:
                    digit = utils.get_digit(text=v)
                    if digit is not None:
                        basic_features.n_rooms = int(digit)
            if "baño" in v:
                digit = utils.get_digit(text=v)
                if digit is not None:
                    basic_features.bathrooms = int(digit)

            if "terraza" in v.lower():
                basic_features.terrace = True
            if "balcón" in v.lower():
                basic_features.balcony = True
            if "Segunda mano" in v:
                basic_features.state = v
            if "Orientación" in v:
                basic_features.orientation = v.replace(
                    "Orientación ",
                    "",
                )
            if "Construido en " in v:
                digit = utils.get_digit(text=v)
                if digit is not None:
                    basic_features.built_in = int(digit)
            if "calefacción" in v.lower():
                basic_features.heating = v
            if "movilidad reducida" in v:
                basic_features.reduced_mobility = v
            if "garaje" in v:
                basic_features.garage = v
            if "Chalet" in v:
                basic_features.semi_detached_house = v

            if "Trastero" in v:
                basic_features.storage_room = True
            if "Aire acondicionado" in v:
                basic_features.air_conditioning = True
            if "Armarios empotrados" in v:
                basic_features.built_in_wardrobes = True
            if "Zonas verdes" in v:
                basic_features.green_areas = True
            if "Piscina" in v:
                basic_features.pool = True
            if "Promoción de obra nueva" in v:
                basic_features.new_building_development = True
        return basic_features

    def to_flat_dict(
        self,
        prefix="basic_features",
    ) -> dict[str, Any]:
        return {f"{prefix}_{k}": v for k, v in asdict(self).items()}


@dataclass
class Building(FlatSerializable, BaseFromRaw):
    elevator: bool | None = None
    inside: bool | None = None
    floor: int | None = None

    @classmethod
    def from_raw(cls, raw: str) -> Self:
        building = cls()
        for v in raw.split(";"):
            if "Sin ascensor" in v:
                building.elevator = False
            elif "Con ascensor" in v:
                building.elevator = True

            if "Entreplanta " in v or "Bajo" in v:
                building.floor = 0
            elif "planta " in v.lower():
                digit = re.search(r"(\d+)(?=ª)", v)
                if digit is not None:
                    building.floor = int(digit.group(1))

            if "interior" in v:
                building.inside = True
            elif "exterior" in v:
                building.inside = False
        return building

    def to_flat_dict(
        self,
        prefix="building",
    ) -> dict[str, Any]:
        return {f"{prefix}_{k}": v for k, v in asdict(self).items()}


@dataclass
class EnergyCerticate(FlatSerializable, BaseFromRaw):
    certificate: str | None = None
    kwh_sqm_year_emissions: float | None = None
    kg_co2_sqm_year_consumption: float | None = None

    @classmethod
    def from_raw(cls, raw: str) -> Self:
        e_cert = cls()
        for v in raw.split(";"):
            if "Emisiones: \n" in v:
                digit = utils.get_digit(text=v)
                if digit is not None:
                    e_cert.kwh_sqm_year_emissions = float(digit)
            if "Consumo: \n" in v:
                digit = utils.get_digit(text=v)
                if digit is not None:
                    e_cert.kg_co2_sqm_year_consumption = float(digit)
            if (
                "Inmueble exento" in v
                or "En trámite" in v
                or "No indicado" in v
            ):
                e_cert.certificate = v
        return e_cert

    def to_flat_dict(
        self,
        prefix="energy_certificate",
    ) -> dict[str, Any]:
        return {f"{prefix}_{k}": v for k, v in asdict(self).items()}


@dataclass
class Status(FlatSerializable, BaseFromRaw):
    tenant: str | None = None

    @classmethod
    def from_raw(cls, raw: str) -> Self:
        status = cls()
        if raw:
            status.tenant = raw
        return status

    def to_flat_dict(self, prefix="status") -> dict[str, Any]:
        return {f"{prefix}_{k}": v for k, v in asdict(self).items()}


@dataclass
class Apartment(FlatSerializable):
    id: int | None = None
    title: str | None = None
    property_price_euros: float | None = None
    price_per_sqm: float | None = None
    monthly_community_fees_euros: float | None = None
    location: str | None = None
    description: str | None = None

    address: Address | None = None

    basic_features: BasicFeatures | None = None
    building: Building | None = None
    energy_certificate: EnergyCerticate | None = None
    status: Status | None = None

    @classmethod
    def from_raw_apartment(cls, raw_apartment: dict) -> Self:
        """Create an Apartment instance from a raw apartment dictionary.

        Args:
            raw_apartment (dict): A dictionary containing raw apartment data.

        Returns:
            Apartment: An instance of the Apartment class populated with the provided data.

        """
        outer_keys_map = {
            "propertyPrice": "property_price_euros",
            "pricePerSqm": "price_per_sqm",
            "communityFees": "monthly_community_fees_euros",
        }

        apartment_info = {}
        for k, value in raw_apartment.items():
            if k in outer_keys_map:
                if value is not None and value != "":
                    match = re.search(
                        r"(\d{1,3}(?:\.\d{3})*)(?:\s?.*)?",
                        value,
                    )
                    if match is not None:
                        apartment_info[outer_keys_map[k]] = float(
                            match.group(1).replace(".", ""),
                        )
            elif k in ["title", "description", "id"]:
                apartment_info[k] = value

        apartment = cls(**apartment_info)

        loaction_v = raw_apartment["location"]
        apartment.location = loaction_v
        apartment.address = Address.from_location(loaction_v)

        apartment.basic_features = BasicFeatures.from_raw(
            raw=raw_apartment["basicFeatures"],
        )

        apartment.building = Building.from_raw(
            raw=raw_apartment["building"],
        )

        apartment.energy_certificate = EnergyCerticate.from_raw(
            raw=raw_apartment["energyCertificate"],
        )

        apartment.status = Status.from_raw(
            raw=raw_apartment["apartmentStatus"],
        )

        return apartment

    def to_flat_dict(self, prefix="apartment") -> dict[str, Any]:
        base = {}

        for field in fields(self):
            value = getattr(self, field.name)
            if (
                isinstance(value, FlatSerializable)
                and value is not None
            ):
                base.update(value.to_flat_dict())
            elif not is_dataclass(value):
                base[f"{prefix}_{field.name}"] = value
        return base
