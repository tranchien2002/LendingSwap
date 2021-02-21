pragma solidity ^0.5.0;
pragma experimental ABIEncoderV2;

interface BandOracleInterface {
    struct ReferenceData {
        uint256 rate;
        uint256 lastUpdatedBase;
        uint256 lastUpdatedQuote;
    }

    function getReferenceData(string calldata _base, string calldata _quote)
        external
        view
        returns (ReferenceData memory);

    function getReferenceDataBulk(
        string[] calldata _bases,
        string[] calldata _quotes
    ) external view returns (ReferenceData[] memory);
}
