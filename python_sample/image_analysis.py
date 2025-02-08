# Copyright (c) Microsoft. All rights reserved.
# Licensed under the MIT license. See LICENSE.md file in the project root for full license information.
"""
Azure AI Vision SDK -- Python Image Analysis Samples

Main entry point for the sample application.
"""

import os
import sys
from azure.core.credentials import AzureKeyCredential
from azure.ai.vision.imageanalysis import ImageAnalysisClient
from azure.ai.vision.imageanalysis.models import VisualFeatures

def print_usage():
    print()
    print(" To run the samples:")
    print()
    print("   python image_analysis.py [--key|-k <your-key>] [--endpoint|-e <your-endpoint>]")
    print()
    print(" Where:")
    print("   <your-key> - A computer vision key you get from your Azure portal.")
    print("     It should be a 32-character HEX number.")
    print("   <your-endpoint> - A computer vision endpoint you get from your Azure portal.")
    print("     It should have the form:")
    print("     https://<your-computer-vision-resource-name>.cognitiveservices.azure.com")
    print()
    print(" As an alternative to specifying the above command line arguments, you can define")
    print(" these environment variables: VISION_KEY and/or VISION_ENDPOINT.")
    print()
    print(" To get this usage help, run:")
    print()
    print("   python image_analysis.py --help|-h")
    print()

def load_secrets():
    key = None
    endpoint = None

    # Check command-line arguments
    for i, arg in enumerate(sys.argv):
        if arg in ["--key", "-k"] and i + 1 < len(sys.argv):
            key = sys.argv[i + 1]
        elif arg in ["--endpoint", "-e"] and i + 1 < len(sys.argv):
            endpoint = sys.argv[i + 1]

    # Check environment variables
    if not key:
        key = os.getenv("VISION_KEY")
    if not endpoint:
        endpoint = os.getenv("VISION_ENDPOINT")

    if not key or not endpoint:
        return None, None

    return key, endpoint

def analyze_image(key, endpoint):
    client = ImageAnalysisClient(endpoint=endpoint, credential=AzureKeyCredential(key))

    image_url1 = "https://aka.ms/azsdk/image-analysis/sample.jpg"  # Replace with your image URL
    image_url2 = "https://encrypted-tbn2.gstatic.com/shopping?q=tbn:ANd9GcR1kz5pVfNQf8xvsfVsO50VVwK7sgonr1IyfUs5p-5p5wHj4FNfGVaBSlnS-yHi6Ab1y2WBHJIMlEDHSvycMh1vv5GEmMw67HQVJCK7Ao2-ZZ1CkUziJJuH388"

    # Analyze the first image for caption, tags, and objects
    visual_features1 = [VisualFeatures.CAPTION, VisualFeatures.TAGS, VisualFeatures.OBJECTS]
    result1 = client.analyze(image_url=image_url1, visual_features=visual_features1)

    if result1.caption:
        print("Caption: {}".format(result1.caption.text))
        print()
        print("----")
        print()
    if result1.tags:
        print("Tags: {}".format(result1.tags))
        print()
        print("----")
        print()
    if result1.objects:
        print("Objects: {}".format(result1.objects))
        print()
        print("----")
        print()

    # Analyze the second image for text (OCR)
    visual_features2 = [VisualFeatures.READ]
    result2 = client.analyze(image_url=image_url2, visual_features=visual_features2)

    if result2.read:
        print("Read - OCR: {}".format(result2.read))
        print()
        print("----")
        print()

if __name__ == "__main__":
    print()
    print(" Azure AI Vision SDK - Image Analysis Samples")
    print()

    if "--help" in sys.argv or "-h" in sys.argv:
        print_usage()
        sys.exit(0)

    key, endpoint = load_secrets()
    if not key or not endpoint:
        print("Error: Missing key or endpoint.")
        print_usage()
        sys.exit(1)

    try:
        analyze_image(key, endpoint)
    except Exception as e:
        print("Error running sample: {}".format(e))

    print()