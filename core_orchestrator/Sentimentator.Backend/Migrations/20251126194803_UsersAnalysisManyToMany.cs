using System;
using Microsoft.EntityFrameworkCore.Migrations;

#nullable disable

namespace Sentimentator.Backend.Migrations
{
    /// <inheritdoc />
    public partial class UsersAnalysisManyToMany : Migration
    {
        /// <inheritdoc />
        protected override void Up(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.DropForeignKey(
                name: "FK_Analyses_AspNetUsers_UserId",
                table: "Analyses");

            migrationBuilder.DropIndex(
                name: "IX_Analyses_UserId",
                table: "Analyses");

            migrationBuilder.DropColumn(
                name: "UserId",
                table: "Analyses");

            migrationBuilder.AddColumn<DateTime>(
                name: "CreatedAt",
                table: "Analyses",
                type: "timestamp with time zone",
                nullable: false,
                defaultValue: new DateTime(1, 1, 1, 0, 0, 0, 0, DateTimeKind.Unspecified));

            migrationBuilder.AddColumn<DateTime>(
                name: "UpdatedAt",
                table: "Analyses",
                type: "timestamp with time zone",
                nullable: false,
                defaultValue: new DateTime(1, 1, 1, 0, 0, 0, 0, DateTimeKind.Unspecified));

            migrationBuilder.CreateTable(
                name: "AnalysisUser",
                columns: table => new
                {
                    AnalysesId = table.Column<Guid>(type: "uuid", nullable: false),
                    UsersId = table.Column<Guid>(type: "uuid", nullable: false)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_AnalysisUser", x => new { x.AnalysesId, x.UsersId });
                    table.ForeignKey(
                        name: "FK_AnalysisUser_Analyses_AnalysesId",
                        column: x => x.AnalysesId,
                        principalTable: "Analyses",
                        principalColumn: "Id",
                        onDelete: ReferentialAction.Cascade);
                    table.ForeignKey(
                        name: "FK_AnalysisUser_AspNetUsers_UsersId",
                        column: x => x.UsersId,
                        principalTable: "AspNetUsers",
                        principalColumn: "Id",
                        onDelete: ReferentialAction.Cascade);
                });

            migrationBuilder.CreateIndex(
                name: "IX_AnalysisUser_UsersId",
                table: "AnalysisUser",
                column: "UsersId");
        }

        /// <inheritdoc />
        protected override void Down(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.DropTable(
                name: "AnalysisUser");

            migrationBuilder.DropColumn(
                name: "CreatedAt",
                table: "Analyses");

            migrationBuilder.DropColumn(
                name: "UpdatedAt",
                table: "Analyses");

            migrationBuilder.AddColumn<Guid>(
                name: "UserId",
                table: "Analyses",
                type: "uuid",
                nullable: true);

            migrationBuilder.CreateIndex(
                name: "IX_Analyses_UserId",
                table: "Analyses",
                column: "UserId");

            migrationBuilder.AddForeignKey(
                name: "FK_Analyses_AspNetUsers_UserId",
                table: "Analyses",
                column: "UserId",
                principalTable: "AspNetUsers",
                principalColumn: "Id");
        }
    }
}
