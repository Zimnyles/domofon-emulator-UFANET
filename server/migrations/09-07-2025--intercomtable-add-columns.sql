ALTER TABLE intercoms
ADD COLUMN CalledApartment INT DEFAULT 0,
ADD COLUMN OpenedDoorApartment INT DEFAULT 0,
ADD COLUMN is_active BOOLEAN DEFAULT TRUE;
