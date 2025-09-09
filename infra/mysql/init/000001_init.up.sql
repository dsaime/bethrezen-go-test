CREATE TABLE `News` (
    `Id` bigint PRIMARY KEY,
    `Title` tinytext NOT NULL,
    `Content` longtext NOT NULL
);

CREATE TABLE `NewsCategories` (
    `NewsId` bigint NOT NULL,
    `CategoryId` bigint NOT NULL,
    PRIMARY KEY (`NewsId`,`CategoryId`)
)