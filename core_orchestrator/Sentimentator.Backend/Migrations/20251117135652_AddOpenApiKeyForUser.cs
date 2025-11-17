using Microsoft.EntityFrameworkCore.Migrations;

#nullable disable

namespace Sentimentator.Backend.Migrations
{
    /// <inheritdoc />
    public partial class AddOpenApiKeyForUser : Migration
    {
        /// <inheritdoc />
        protected override void Up(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.AddColumn<string>(
                name: "OpenAiApiKey",
                table: "AspNetUsers",
                type: "text",
                nullable: true);
        }

        /// <inheritdoc />
        protected override void Down(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.DropColumn(
                name: "OpenAiApiKey",
                table: "AspNetUsers");
        }
    }
}
