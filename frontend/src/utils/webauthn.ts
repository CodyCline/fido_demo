// Base64 to ArrayBuffer
export function bufferDecode(value:any) {
    return Uint8Array.from(atob(value), c => c.charCodeAt(0));
}

// ArrayBuffer to URLBase64
export function bufferEncode(value:any) {
    const arr: any = new Uint8Array(value)
    return btoa(String.fromCharCode.apply(null, arr))
        .replace(/\+/g, "-")
        .replace(/\//g, "_")
        .replace(/=/g, "");
}