export const shortAddres = async address => {
  return address && address.length > 0
    ? address.substring(0, 6) + '...' + address.substring(address.length - 5, address.length - 1)
    : address;
};
