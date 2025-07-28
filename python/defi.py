import requests
import datetime

# === CONFIGURATION ===
NFT_ID = "3875527"  # Replace with your Uniswap V3 Position NFT ID
GRAPH_URL = "https://gateway.thegraph.com/api/9e977614ac5ec3f609cf47d3bffa1174/subgraphs/id/FQ6JYszEKApsBpAmiHesRsd9Ygc6mzmpNRANeVQFYoVX"  # Uniswap V3 on Arbitrum

# === QUERY GRAPH ===
def fetch_uniswap_position_fees(nft_id):
    query = """
    {
      position(id: \"""" + nft_id + """\") {
        id
        collectedFeesToken0
        collectedFeesToken1
        pool {
          token0 {
            symbol
            decimals
          }
          token1 {
            symbol
            decimals
          }
        }
        createdAtTimestamp
      }
    }
    """
    response = requests.post(GRAPH_URL, json={"query": query})
    
    if response.status_code == 200:
        data = response.json().get("data", {}).get("position")
        return data
    else:
        print("Error fetching data:", response.text)
        return None

# === CALCULATE AVG DAILY FEE ===
def calculate_avg_daily_fee(nft_id):
    data = fetch_uniswap_position_fees(nft_id)
    
    if not data:
        print("No data found for NFT ID:", nft_id)
        return None

    # Extract collected fees and token decimals
    fees_token0 = float(data["collectedFeesToken0"]) / (10 ** int(data["pool"]["token0"]["decimals"]))
    fees_token1 = float(data["collectedFeesToken1"]) / (10 ** int(data["pool"]["token1"]["decimals"]))
    
    # Extract creation date
    created_at = int(data["createdAtTimestamp"])
    created_date = datetime.datetime.utcfromtimestamp(created_at)
    today_date = datetime.datetime.utcnow()
    days_active = (today_date - created_date).days or 1  # Avoid division by zero
    
    # Compute average daily fee
    avg_fee_token0 = fees_token0 / days_active
    avg_fee_token1 = fees_token1 / days_active

    # Display results
    print(f"🔹 Position NFT ID: {nft_id}")
    print(f"📅 Days Active: {days_active}")
    print(f"💰 Total Collected Fees: {fees_token0:.6f} {data['pool']['token0']['symbol']} + {fees_token1:.6f} {data['pool']['token1']['symbol']}")
    print(f"📊 Avg Daily Fee: {avg_fee_token0:.6f} {data['pool']['token0']['symbol']} + {avg_fee_token1:.6f} {data['pool']['token1']['symbol']}")
    
    return avg_fee_token0, avg_fee_token1

# === RUN SCRIPT ===
calculate_avg_daily_fee(NFT_ID)
