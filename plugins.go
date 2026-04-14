package main

// Plugin holds the display name and Windows zip download URL for one plugin.
// Source: https://www.patreon.com/posts/all-download-34851999
type Plugin struct {
	Name       string // display name shown in UI
	ZipURL     string
	BundleName string // actual .vst3 bundle name inside the zip; defaults to Name if empty
	PageURL    string // plugin page on analogobsession.com
	Desc       string // one-line description for tooltip / info panel
}

// Bundle returns the name used for .vst3 file operations.
func (p Plugin) Bundle() string {
	if p.BundleName != "" {
		return p.BundleName
	}
	return p.Name
}

// Plugins is the complete list of Analog Obsession plugins that ship a Windows zip.
// Sorted ascending alphabetically (A → Z), case-insensitive.
// Bundles (COLOR_BUNDLE, FET BUNDLE, etc.) omitted — Patreon-post-only links.
var Plugins = []Plugin{
	{
		Name:    "AHEAD",
		ZipURL:  "https://analogobsession.com/wp-content/uploads/2024/07/AHEAD_1.0.zip",
		PageURL: "https://analogobsession.com/color-preamp-saturation/ahead/",
		Desc:    "Four-model guitar amp with three-band tone stack",
	},
	{
		Name:    "ATONE",
		ZipURL:  "https://analogobsession.com/wp-content/uploads/2024/04/ATONE_1.0.zip",
		PageURL: "https://analogobsession.com/channel-strip/atone/",
		Desc:    "Altec-style 3-band EQ + 436 compressor channel strip",
	},
	{
		Name:    "ATTRACTOR",
		ZipURL:  "https://analogobsession.com/wp-content/uploads/2025/02/ATTRACTOR_1.0.zip",
		PageURL: "https://analogobsession.com/dynamicprocessing/attractor/",
		Desc:    "Transient processor with independent attack/release paths",
	},
	{
		Name:    "BLENDEQ",
		ZipURL:  "https://analogobsession.com/wp-content/uploads/2023/05/BLENDEQ_2.0.zip",
		PageURL: "https://analogobsession.com/equalization/blendeq/",
		Desc:    "Four-band EQ blending American and British circuit models",
	},
	{
		Name:    "BlackVibe",
		ZipURL:  "https://analogobsession.com/wp-content/uploads/2022/08/BlackVibe_2.0.zip",
		PageURL: "https://analogobsession.com/color-preamp-saturation/blackvibe/",
		Desc:    "Fender-style clean tube amp with tremolo, no cabinet sim",
	},
	{
		Name:    "BritChannel",
		ZipURL:  "https://analogobsession.com/wp-content/uploads/2023/03/BritChannel_7.0.zip",
		PageURL: "https://analogobsession.com/channel-strip/britchannel/",
		Desc:    "British Type 73 mic/line preamp with 3-band EQ",
	},
	{
		Name:    "Britpressor",
		ZipURL:  "https://analogobsession.com/wp-content/uploads/2023/03/Britpressor_3.0.zip",
		PageURL: "https://analogobsession.com/dynamicprocessing/britpressor/",
		Desc:    "Vintage solid-state compressor with multiband sidechain EQ",
	},
	{
		Name:    "BritPRE",
		ZipURL:  "https://analogobsession.com/wp-content/uploads/2023/03/BritPre_2.0.zip",
		PageURL: "https://analogobsession.com/color-preamp-saturation/britpre/",
		Desc:    "British mic/line preamp with variable HP/LP filters",
	},
	{
		Name:    "BUSTERse",
		ZipURL:  "https://analogobsession.com/wp-content/uploads/2022/12/BUSTERse_7.0.zip",
		PageURL: "https://analogobsession.com/dynamicprocessing/busterse/",
		Desc:    "Console bus compressor with filter and transient sidechain",
	},
	{
		Name:    "BXQ",
		ZipURL:  "https://analogobsession.com/wp-content/uploads/2023/05/BXQ_5.0.zip",
		PageURL: "https://analogobsession.com/equalization/bxq/",
		Desc:    "Mastering Baxandall EQ with L/R and M/S processing",
	},
	{
		Name:    "CHANNEV",
		ZipURL:  "https://analogobsession.com/wp-content/uploads/2023/02/CHANNEV_2.0.zip",
		PageURL: "https://analogobsession.com/channel-strip/channev/",
		Desc:    "Full channel strip: preamp, de-esser, EQ, comp, limiter, tape",
	},
	{
		Name:    "Chopa",
		ZipURL:  "https://analogobsession.com/wp-content/uploads/2022/06/Chopa_3.0.zip",
		PageURL: "https://analogobsession.com/modulation/chopa/",
		Desc:    "Tremolo/chopper effect with adjustable rate and shape",
	},
	{
		Name:    "CITE",
		ZipURL:  "https://analogobsession.com/wp-content/uploads/2024/05/CITE_1.0.zip",
		PageURL: "https://analogobsession.com/color-preamp-saturation/cite/",
		Desc:    "High-frequency air processor with frequency area selection",
	},
	{
		Name:    "COMBOX",
		ZipURL:  "https://analogobsession.com/wp-content/uploads/2023/09/COMBOX_6.0.zip",
		PageURL: "https://analogobsession.com/equalization/combox/",
		Desc:    "American 3-band inductor EQ with variable HP/LP filters",
	},
	{
		Name:    "COMPER",
		ZipURL:  "https://analogobsession.com/wp-content/uploads/2022/11/COMPER_1.1.zip",
		PageURL: "https://analogobsession.com/dynamicprocessing/comper/",
		Desc:    "Multi-mode serial compressor: VCA, FET, and Opto circuits",
	},
	{
		Name:    "dBComp",
		ZipURL:  "https://analogobsession.com/wp-content/uploads/2022/12/dBComp_2.0.zip",
		PageURL: "https://analogobsession.com/dynamicprocessing/dbcomp/",
		Desc:    "Vintage-inspired analog compressor",
	},
	{
		Name:    "Distox",
		ZipURL:  "https://analogobsession.com/wp-content/uploads/2024/02/CB_Distox_1.0.zip",
		PageURL: "https://analogobsession.com/color-preamp-saturation/distox/",
		Desc:    "Overdrive/distortion with multiple circuit models",
	},
	{
		Name:    "DrGate",
		ZipURL:  "https://analogobsession.com/wp-content/uploads/2024/11/DrGate_1.0.zip",
		PageURL: "https://analogobsession.com/dynamicprocessing/drgate/",
		Desc:    "Noise gate with precise threshold and release control",
	},
	{
		Name:    "EDComp",
		ZipURL:  "https://analogobsession.com/wp-content/uploads/2025/05/EDComp_1.0.zip",
		PageURL: "https://analogobsession.com/dynamicprocessing/edcomp/",
		Desc:    "Vintage electro-dynamic compressor emulation",
	},
	{
		Name:    "FetCB",
		ZipURL:  "https://analogobsession.com/wp-content/uploads/2023/11/FetCB_1.0.zip",
		PageURL: "https://analogobsession.com/dynamicprocessing/fetcb/",
		Desc:    "FET-style compressor with classic solid-state character",
	},
	{
		Name:    "FETish",
		ZipURL:  "https://analogobsession.com/wp-content/uploads/2023/04/FETish_6.0.zip",
		PageURL: "https://analogobsession.com/dynamicprocessing/fetish/",
		Desc:    "FET compressor emulation based on classic hardware limiters",
	},
	{
		Name:    "FetDrive",
		ZipURL:  "https://analogobsession.com/wp-content/uploads/2023/04/FetDrive_3.0.zip",
		PageURL: "https://analogobsession.com/color-preamp-saturation/fetdrive/",
		Desc:    "FET-based drive and saturation processor",
	},
	{
		Name:    "FetSnap",
		ZipURL:  "https://analogobsession.com/wp-content/uploads/2023/04/FetSnap_3.0.zip",
		PageURL: "https://analogobsession.com/dynamicprocessing/fetsnap/",
		Desc:    "FET transient shaper with snap and sustain control",
	},
	{
		Name:    "FIVER",
		ZipURL:  "https://analogobsession.com/wp-content/uploads/2026/01/FIVER_5.0.zip",
		PageURL: "https://analogobsession.com/equalization/fiver/",
		Desc:    "Five-band vintage-style parametric equalizer",
	},
	{
		Name:    "Frank",
		ZipURL:  "https://analogobsession.com/wp-content/uploads/2023/05/Frank_3.0.zip",
		PageURL: "https://analogobsession.com/equalization/frank/",
		Desc:    "Classic four-band analog equalizer",
	},
	{
		Name:    "FrankCS",
		ZipURL:  "https://analogobsession.com/wp-content/uploads/2023/05/FrankCS_2.0.zip",
		PageURL: "https://analogobsession.com/channel-strip/frankcs/",
		Desc:    "Frank EQ-based channel strip with compression",
	},
	{
		Name:       "G395a",
		ZipURL:     "https://analogobsession.com/wp-content/uploads/2023/05/G395a_2.0.zip",
		BundleName: "GGGGa",
		PageURL:    "https://analogobsession.com/equalization/g395a/",
		Desc:       "Neve-inspired 1073-style mic preamp EQ",
	},
	{
		Name:    "GrapHack",
		ZipURL:  "https://analogobsession.com/wp-content/uploads/2023/04/GrapHack_1.1.zip",
		PageURL: "https://analogobsession.com/equalization/graphack/",
		Desc:    "Graphic equalizer with hacksaw-style frequency shaping",
	},
	{
		Name:    "Harqules",
		ZipURL:  "https://analogobsession.com/wp-content/uploads/2023/05/Harqules_5.0.zip",
		PageURL: "https://analogobsession.com/equalization/harqules/",
		Desc:    "Vintage-style program equalizer",
	},
	{
		Name:    "HLQSE",
		ZipURL:  "https://analogobsession.com/wp-content/uploads/2023/05/HLQSE_5.0.zip",
		PageURL: "https://analogobsession.com/equalization/hlqse/",
		Desc:    "Type 69-style inductor program equalizer",
	},
	{
		Name:    "INDEQ",
		ZipURL:  "https://analogobsession.com/wp-content/uploads/2025/08/INDEQ_1.0.zip",
		PageURL: "https://analogobsession.com/equalization/indeq/",
		Desc:    "Pure inductor three-band EQ with highpass filter",
	},
	{
		Name:    "KABIN",
		ZipURL:  "https://analogobsession.com/wp-content/uploads/2024/08/KABIN_3.0.zip",
		PageURL: "https://analogobsession.com/color-preamp-saturation/kabin/",
		Desc:    "Guitar cabinet modeler with speaker size and mic placement",
	},
	{
		Name:    "KOLIN",
		ZipURL:  "https://analogobsession.com/wp-content/uploads/2023/02/Kolin_5.1.zip",
		PageURL: "https://analogobsession.com/dynamicprocessing/kolin/",
		Desc:    "Vari-Mu tube limiting amplifier",
	},
	{
		Name:    "KolinMB",
		ZipURL:  "https://analogobsession.com/wp-content/uploads/2023/06/KolinMB_1.0.zip",
		PageURL: "https://analogobsession.com/dynamicprocessing/kolinmb/",
		Desc:    "Multiband Vari-Mu tube limiting amplifier",
	},
	{
		Name:    "KONSOL",
		ZipURL:  "https://analogobsession.com/wp-content/uploads/2021/02/KONSOL_3.0.zip",
		PageURL: "https://analogobsession.com/channel-strip/konsol/",
		Desc:    "Multi-model console strip: tube, transistor, and op-amp modes",
	},
	{
		Name:    "LALA",
		ZipURL:  "https://analogobsession.com/wp-content/uploads/2023/01/LALA_3.1.zip",
		PageURL: "https://analogobsession.com/dynamicprocessing/lala/",
		Desc:    "LA-2A optical limiter/compressor with sidechain filter",
	},
	{
		Name:    "LAEA",
		ZipURL:  "https://analogobsession.com/wp-content/uploads/2026/03/LAEA_1.2.zip",
		PageURL: "https://analogobsession.com/dynamicprocessing/laea/",
		Desc:    "UREI LA-3A two-knob optical compressor emulation",
	},
	{
		Name:    "LOADES",
		ZipURL:  "https://analogobsession.com/wp-content/uploads/2024/09/LOADES_2.0.zip",
		PageURL: "https://analogobsession.com/dynamicprocessing/loades/",
		Desc:    "De-esser for controlling sibilance in vocals and instruments",
	},
	{
		Name:    "LOVEND",
		ZipURL:  "https://analogobsession.com/wp-content/uploads/2021/05/LOVEND_2.0.zip",
		PageURL: "https://analogobsession.com/color-preamp-saturation/lovend/",
		Desc:    "Harmonic bass enhancer adding warmth to low frequencies",
	},
	{
		Name:    "MAXBAX",
		ZipURL:  "https://analogobsession.com/wp-content/uploads/2023/12/MAXBAX_2.0.zip",
		PageURL: "https://analogobsession.com/equalization/maxbax/",
		Desc:    "Passive Baxandall EQ with stepped boost/cut per band",
	},
	{
		Name:    "MERICA",
		ZipURL:  "https://analogobsession.com/wp-content/uploads/2023/09/MERICA_4.0.zip",
		PageURL: "https://analogobsession.com/equalization/merica/",
		Desc:    "American console proportional-Q three-band equalizer",
	},
	{
		Name:    "MidBoss",
		ZipURL:  "https://analogobsession.com/wp-content/uploads/2025/10/MidBoss_1.0.zip",
		PageURL: "https://analogobsession.com/dynamicprocessing/midboss/",
		Desc:    "Mid-frequency saturator with six saturation type options",
	},
	{
		Name:    "MPReq",
		ZipURL:  "https://analogobsession.com/wp-content/uploads/2023/12/MPReq_5.0.zip",
		PageURL: "https://analogobsession.com/equalization/mpreq/",
		Desc:    "Vintage mic preamp with two-band program equalizer",
	},
	{
		Name:    "MythPre",
		ZipURL:  "https://analogobsession.com/wp-content/uploads/2025/04/MythPre_1.0.zip",
		PageURL: "https://analogobsession.com/color-preamp-saturation/mythpre/",
		Desc:    "Tube mic preamp with tone shaper and low-end emphasis",
	},
	{
		Name:       "N492ME",
		ZipURL:     "https://analogobsession.com/wp-content/uploads/2023/12/N492ME_7.0.zip",
		BundleName: "NEQME",
		PageURL:    "https://analogobsession.com/equalization/n492me/",
		Desc:       "Neumann-style 4-band parametric EQ with M/S and L/R modes",
	},
	{
		Name:    "OAQ",
		ZipURL:  "https://analogobsession.com/wp-content/uploads/2023/12/OAQ_5.0.zip",
		PageURL: "https://analogobsession.com/equalization/oaq/",
		Desc:    "Mastering EQ with dual-mono/stereo/M-S modes and Drive knob",
	},
	{
		Name:    "OSS",
		ZipURL:  "https://analogobsession.com/wp-content/uploads/2022/12/OSS_6.0.zip",
		PageURL: "https://analogobsession.com/dynamicprocessing/oss/",
		Desc:    "Analog compressor with transient manipulation controls",
	},
	{
		Name:    "PEDALz",
		ZipURL:  "https://analogobsession.com/wp-content/uploads/2024/08/PEDALz_2.0.zip",
		PageURL: "https://analogobsession.com/color-preamp-saturation/pedalz/",
		Desc:    "Five-in-one drive pedal with multiple distortion models",
	},
	{
		Name:    "POORTEC",
		ZipURL:  "https://analogobsession.com/wp-content/uploads/2025/01/POORTEC_1.0.zip",
		PageURL: "https://analogobsession.com/equalization/poortec/",
		Desc:    "Pultec EQP-1A-inspired passive program equalizer",
	},
	{
		Name:    "PreBOX",
		ZipURL:  "https://analogobsession.com/wp-content/uploads/2024/02/CB_PreBOX_1.0.zip",
		PageURL: "https://analogobsession.com/color-preamp-saturation/prebox/",
		Desc:    "Multi-model microphone preamp simulator",
	},
	{
		Name:    "RazorClip",
		ZipURL:  "https://analogobsession.com/wp-content/uploads/2026/02/RazorClip_1.0.zip",
		PageURL: "https://analogobsession.com/color-preamp-saturation/razorclip/",
		Desc:    "Analog clipper with five different clipping circuit models",
	},
	{
		Name:    "Rare",
		ZipURL:  "https://analogobsession.com/wp-content/uploads/2023/10/RB_Rare_1.0.zip",
		PageURL: "https://analogobsession.com/equalization/rare/",
		Desc:    "Vintage program EQ with L/R and M/S processing",
	},
	{
		Name:    "RareSE",
		ZipURL:  "https://analogobsession.com/wp-content/uploads/2023/10/RB_RareSE_1.0.zip",
		PageURL: "https://analogobsession.com/equalization/rarese/",
		Desc:    "Vintage program EQ with extended processing capabilities",
	},
	{
		Name:    "ReLife",
		ZipURL:  "https://analogobsession.com/wp-content/uploads/2024/10/ReLife_2.0.zip",
		PageURL: "https://analogobsession.com/dynamicprocessing/relife/",
		Desc:    "Dynamic enhancer adding life and presence to flat tracks",
	},
	{
		Name:       "Room041",
		ZipURL:     "https://analogobsession.com/wp-content/uploads/2024/09/Room041_2.0.zip",
		BundleName: "RevO",
		PageURL:    "https://analogobsession.com/color-preamp-saturation/room041/",
		Desc:       "Room and plate reverb with natural characteristics",
	},
	{
		Name:    "SPECOMP",
		ZipURL:  "https://analogobsession.com/wp-content/uploads/2022/12/SPECOMP_2.0.zip",
		PageURL: "https://analogobsession.com/dynamicprocessing/specomp/",
		Desc:    "Spectral analog compressor managing the full frequency range",
	},
	{
		Name:    "SSQ",
		ZipURL:  "https://analogobsession.com/wp-content/uploads/2023/12/SSQ_7.0.zip",
		PageURL: "https://analogobsession.com/equalization/ssq/",
		Desc:    "Famous console-style four-band equalizer",
	},
	{
		Name:    "STEQ",
		ZipURL:  "https://analogobsession.com/wp-content/uploads/2023/12/STEQ_4.0.zip",
		PageURL: "https://analogobsession.com/equalization/steq/",
		Desc:    "Small format console EQ",
	},
	{
		Name:    "SweetDrums",
		ZipURL:  "https://analogobsession.com/wp-content/uploads/2021/11/SweetDrums_4.0.zip",
		PageURL: "https://analogobsession.com/dynamicprocessing/sweetdrums/",
		Desc:    "One-knob drum processor with EQ, compression, and saturation",
	},
	{
		Name:    "SweetVox",
		ZipURL:  "https://analogobsession.com/wp-content/uploads/2021/11/SweetVox_4.0.zip",
		PageURL: "https://analogobsession.com/dynamicprocessing/sweetvox/",
		Desc:    "Intelligent vocal processor with built-in de-esser",
	},
	{
		Name:    "TILTA",
		ZipURL:  "https://analogobsession.com/wp-content/uploads/2023/12/TILTA_2.0.zip",
		PageURL: "https://analogobsession.com/equalization/tilta/",
		Desc:    "Tilt equalizer for broad tonal balance adjustment",
	},
	{
		Name:    "TRAX",
		ZipURL:  "https://analogobsession.com/wp-content/uploads/2024/02/TRAX_2.0.zip",
		PageURL: "https://analogobsession.com/dynamicprocessing/trax/",
		Desc:    "Transient and dynamics processor",
	},
	{
		Name:    "TREQ",
		ZipURL:  "https://analogobsession.com/wp-content/uploads/2023/07/TREQ_6.0.zip",
		PageURL: "https://analogobsession.com/equalization/treq/",
		Desc:    "T-style classic four-band equalizer",
	},
	{
		Name:    "TUBA",
		ZipURL:  "https://analogobsession.com/wp-content/uploads/2023/02/TUBA_3.0.zip",
		PageURL: "https://analogobsession.com/color-preamp-saturation/tuba/",
		Desc:    "Tube console preamp with two-band EQ and mic/line modes",
	},
	{
		Name:    "TuPRE",
		ZipURL:  "https://analogobsession.com/wp-content/uploads/2023/02/TuPRE_3.0.zip",
		PageURL: "https://analogobsession.com/color-preamp-saturation/tupre/",
		Desc:    "Tube line amplifier with custom passive program equalizer",
	},
	{
		Name:    "UREQ",
		ZipURL:  "https://analogobsession.com/wp-content/uploads/2023/08/UREQ_1.0.zip",
		PageURL: "https://analogobsession.com/equalization/ureq/",
		Desc:    "U-style classic two-band program equalizer",
	},
	{
		Name:    "VariMoon",
		ZipURL:  "https://analogobsession.com/wp-content/uploads/2022/12/VariMoon_6.0.zip",
		PageURL: "https://analogobsession.com/dynamicprocessing/varimoon/",
		Desc:    "Variable-mu tube compressor with classic limiting character",
	},
	{
		Name:    "YALA",
		ZipURL:  "https://analogobsession.com/wp-content/uploads/2022/12/YALA_6.0.zip",
		PageURL: "https://analogobsession.com/dynamicprocessing/yala/",
		Desc:    "Iconic Vari-Mu limiting amplifier with extra features",
	},
}
