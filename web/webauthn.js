var webauthn;
/******/ (() => { // webpackBootstrap
/******/ 	"use strict";
/******/ 	// The require scope
/******/ 	var __webpack_require__ = {};
/******/ 	
/************************************************************************/
/******/ 	/* webpack/runtime/define property getters */
/******/ 	(() => {
/******/ 		// define getter functions for harmony exports
/******/ 		__webpack_require__.d = (exports, definition) => {
/******/ 			for(var key in definition) {
/******/ 				if(__webpack_require__.o(definition, key) && !__webpack_require__.o(exports, key)) {
/******/ 					Object.defineProperty(exports, key, { enumerable: true, get: definition[key] });
/******/ 				}
/******/ 			}
/******/ 		};
/******/ 	})();
/******/ 	
/******/ 	/* webpack/runtime/hasOwnProperty shorthand */
/******/ 	(() => {
/******/ 		__webpack_require__.o = (obj, prop) => (Object.prototype.hasOwnProperty.call(obj, prop))
/******/ 	})();
/******/ 	
/******/ 	/* webpack/runtime/make namespace object */
/******/ 	(() => {
/******/ 		// define __esModule on exports
/******/ 		__webpack_require__.r = (exports) => {
/******/ 			if(typeof Symbol !== 'undefined' && Symbol.toStringTag) {
/******/ 				Object.defineProperty(exports, Symbol.toStringTag, { value: 'Module' });
/******/ 			}
/******/ 			Object.defineProperty(exports, '__esModule', { value: true });
/******/ 		};
/******/ 	})();
/******/ 	
/************************************************************************/
var __webpack_exports__ = {};
// ESM COMPAT FLAG
__webpack_require__.r(__webpack_exports__);

// EXPORTS
__webpack_require__.d(__webpack_exports__, {
  Base64URL: () => (/* reexport */ Base64URL),
  WebAuthnClient: () => (/* reexport */ WebAuthnClient)
});

;// CONCATENATED MODULE: ./src/codec_base64url.ts
function encodeBase64URL(buffer) {
    var binary = '';
    var bytes = new Uint8Array(buffer);
    for (var i = 0; i < bytes.byteLength; i++) {
        binary += String.fromCharCode(bytes[i]);
    }
    return btoa(binary)
        .replace(/\+/g, '-')
        .replace(/\//g, '_')
        .replace(/=/g, '');
}
function decodeBase64URL(str) {
    return Uint8Array.from(atob(str.replace(/-/g, '+').replace(/_/g, '/')), function (c) { return c.charCodeAt(0); }).buffer;
}
var Base64URL = {
    encode: encodeBase64URL,
    decode: decodeBase64URL,
};

;// CONCATENATED MODULE: ./src/types_authentication.ts
function convertAuthenticationChallenge(input, codec) {
    return {
        rpId: input.rpId,
        // userVerification: 'preferred',
        challenge: codec.decode(input.challenge),
        allowCredentials: input.allowCredentials.map(function (cred) { return ({
            id: codec.decode(cred.id),
            type: cred.type,
            // transports: ['usb', 'ble', 'nfc'],
        }); }),
    };
}
function convertAuthenticationResponse(cred, token, challenge, codec) {
    var response = cred.response;
    return {
        token: token,
        challenge: challenge,
        credentialId: codec.encode(cred.rawId),
        response: {
            authenticatorData: codec.encode(response.authenticatorData),
            clientDataJSON: codec.encode(response.clientDataJSON),
            signature: codec.encode(response.signature),
            userHandle: response.userHandle
                ? codec.encode(response.userHandle)
                : null,
        },
    };
}

;// CONCATENATED MODULE: ./src/types_registration.ts
function convertRegistrationChallenge(input, codec) {
    return {
        rp: input.rp,
        user: {
            id: new Uint8Array(input.user.id.split('').map(function (c) { return c.charCodeAt(0); })),
            name: input.user.name,
            displayName: input.user.displayName,
        },
        // excludeCredentials: [],
        pubKeyCredParams: input.pubKeyCredParams,
        attestation: 'direct',
        // authenticatorSelection: {
        // 	authenticatorAttachment: 'cross-platform',
        // 	userVerification: 'preferred',
        // 	requireResidentKey: false,
        // },
        challenge: codec.decode(input.challenge),
    };
}
function convertRegistrationResponse(cred, token, challenge, codec) {
    var response = cred.response;
    return {
        token: token,
        challenge: challenge,
        credentialId: codec.encode(cred.rawId),
        response: {
            clientDataJSON: codec.encode(response.clientDataJSON),
            attestationObject: codec.encode(response.attestationObject),
        },
    };
}

;// CONCATENATED MODULE: ./src/webauthn.ts
var __awaiter = (undefined && undefined.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
var __generator = (undefined && undefined.__generator) || function (thisArg, body) {
    var _ = { label: 0, sent: function() { if (t[0] & 1) throw t[1]; return t[1]; }, trys: [], ops: [] }, f, y, t, g;
    return g = { next: verb(0), "throw": verb(1), "return": verb(2) }, typeof Symbol === "function" && (g[Symbol.iterator] = function() { return this; }), g;
    function verb(n) { return function (v) { return step([n, v]); }; }
    function step(op) {
        if (f) throw new TypeError("Generator is already executing.");
        while (g && (g = 0, op[0] && (_ = 0)), _) try {
            if (f = 1, y && (t = op[0] & 2 ? y["return"] : op[0] ? y["throw"] || ((t = y["return"]) && t.call(y), 0) : y.next) && !(t = t.call(y, op[1])).done) return t;
            if (y = 0, t) op = [op[0] & 2, t.value];
            switch (op[0]) {
                case 0: case 1: t = op; break;
                case 4: _.label++; return { value: op[1], done: false };
                case 5: _.label++; y = op[1]; op = [0]; continue;
                case 7: op = _.ops.pop(); _.trys.pop(); continue;
                default:
                    if (!(t = _.trys, t = t.length > 0 && t[t.length - 1]) && (op[0] === 6 || op[0] === 2)) { _ = 0; continue; }
                    if (op[0] === 3 && (!t || (op[1] > t[0] && op[1] < t[3]))) { _.label = op[1]; break; }
                    if (op[0] === 6 && _.label < t[1]) { _.label = t[1]; t = op; break; }
                    if (t && _.label < t[2]) { _.label = t[2]; _.ops.push(op); break; }
                    if (t[2]) _.ops.pop();
                    _.trys.pop(); continue;
            }
            op = body.call(thisArg, _);
        } catch (e) { op = [6, e]; y = 0; } finally { f = t = 0; }
        if (op[0] & 5) throw op[1]; return { value: op[0] ? op[1] : void 0, done: true };
    }
};



var WebAuthnClient = /** @class */ (function () {
    function WebAuthnClient(codec) {
        if (codec === void 0) { codec = Base64URL; }
        this.codec = codec;
    }
    WebAuthnClient.prototype.register = function (challenge, timeout) {
        return __awaiter(this, void 0, void 0, function () {
            var options, cred;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        options = convertRegistrationChallenge(challenge, this.codec);
                        options.timeout = timeout;
                        return [4 /*yield*/, navigator.credentials.create({ publicKey: options })];
                    case 1:
                        cred = _a.sent();
                        if (!cred || !(cred instanceof PublicKeyCredential))
                            throw new Error('invalid credential');
                        // Convert the credential into a RegistrationResponse
                        return [2 /*return*/, convertRegistrationResponse(cred, challenge.token, challenge.challenge, this.codec)];
                }
            });
        });
    };
    WebAuthnClient.prototype.authenticate = function (challenge, timeout) {
        return __awaiter(this, void 0, void 0, function () {
            var options, assertion;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        options = convertAuthenticationChallenge(challenge, this.codec);
                        options.timeout = timeout;
                        return [4 /*yield*/, navigator.credentials.get({
                                publicKey: options,
                            })];
                    case 1:
                        assertion = _a.sent();
                        if (!assertion || !(assertion instanceof PublicKeyCredential))
                            throw new Error('invalid credential assertion');
                        // Convert the credential into an AuthenticationResponse
                        return [2 /*return*/, convertAuthenticationResponse(assertion, challenge.token, challenge.challenge, this.codec)];
                }
            });
        });
    };
    return WebAuthnClient;
}());


;// CONCATENATED MODULE: ./src/index.ts




webauthn = __webpack_exports__;
/******/ })()
;