function getObjectID(s) {
    return s.toLowerCase().replaceAll("-", " ")
        .replaceAll(".", " ")
        .replaceAll("/", " ")
        .replaceAll(" ", "");
}